package checks

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/helm/chart-testing/v3/pkg/chart"
	"github.com/helm/chart-testing/v3/pkg/config"
	"github.com/helm/chart-testing/v3/pkg/exec"
	"github.com/helm/chart-testing/v3/pkg/util"
	"github.com/imdario/mergo"
	"github.com/redhat-certification/chart-verifier/pkg/tool"
	"gopkg.in/yaml.v3"
)

// buildChartTestingConfiguration computes the chart testing related
// configuration from the given check options.
func buildChartTestingConfiguration(opts *CheckOptions) config.Configuration {

	// cfg will be populated with options gathered from the input
	// check options.
	cfg := config.Configuration{
		BuildId:           opts.ViperConfig.GetString("buildId"),
		Upgrade:           opts.ViperConfig.GetBool("upgrade"),
		SkipMissingValues: opts.ViperConfig.GetBool("skipMissingValues"),
		ReleaseLabel:      opts.ViperConfig.GetString("releaseLabel"),
		Namespace:         opts.ViperConfig.GetString("namespace"),
		HelmExtraArgs:     opts.ViperConfig.GetString("helmExtraArgs"),
	}

	if len(cfg.BuildId) == 0 {
		cfg.BuildId = "build-" + util.RandomString(6)
	}

	if len(cfg.ReleaseLabel) == 0 {
		cfg.ReleaseLabel = "app.kubernetes.io/instance"
	}

	if len(cfg.Namespace) == 0 {
		// Namespace() returns "default" unless has been overriden
		// through environment variables.
		cfg.Namespace = opts.HelmEnvSettings.Namespace()
	}

	return cfg
}

// ChartTesting partially integrates the chart-testing project in chart-verifier.
//
// Unfortunately it wasn't easy as initially expect to integrate
// chart-testing as a lib in the project, including the main
// orchestration logic. The ChartTesting function is the
// interpretation the main logic chart-testing carries, and other
// functions used in this context were also ported from
// chart-verifier.
//
// Helm and kubectl are requirements in the system executing the check
// in order to orchestrate the install, upgrade and chart testing
// phases.
func ChartTesting(opts *CheckOptions) (Result, error) {

	cfg := buildChartTestingConfiguration(opts)
	procExec := exec.NewProcessExecutor(cfg.Debug)
	extraArgs := strings.Fields(cfg.HelmExtraArgs)
	helm := tool.NewHelm(procExec, extraArgs)
	kubectl := tool.NewKubectl(procExec)

	_, path, err := LoadChartFromURI(opts.URI)
	if err != nil {
		return NewResult(false, err.Error()), nil
	}

	chrt, err := chart.NewChart(path)
	if err != nil {
		return NewResult(false, err.Error()), nil
	}

	if cfg.Upgrade {
		oldChrt, err := getChartPreviousVersion(chrt)
		if err != nil {
			return NewResult(
					false,
					fmt.Sprintf("skipping upgrade test of '%s' because no previous chart is available", chrt.Yaml().Name)),
				nil
		}
		breakingChangeAllowed, err := util.BreakingChangeAllowed(oldChrt.Yaml().Version, chrt.Yaml().Version)
		if !breakingChangeAllowed {
			return NewResult(
					false,
					fmt.Sprintf("Skipping upgrade test of '%s' because breaking changes are not allowed for chart", chrt)),
				nil
		} else if err != nil {
			return NewResult(false, err.Error()), nil
		}
		result := upgradeAndTestChart(cfg, oldChrt, chrt, helm, kubectl)

		if result.Error != nil {
			return NewResult(false, result.Error.Error()), nil
		}
	} else {
		result := installAndTestChartRelease(cfg, chrt, helm, kubectl, opts.Values)
		if result.Error != nil {
			return NewResult(false, result.Error.Error()), nil
		}
	}

	return NewResult(true, ChartTestingSuccess), nil
}

// generateInstallConfig extracts required information to install a
// release and builds a clenup function to be used after tests are
// executed.
func generateInstallConfig(
	cfg config.Configuration,
	chrt *chart.Chart,
	helm tool.Helm,
	kubectl tool.Kubectl,
) (namespace, release, releaseSelector string, cleanup func()) {
	if cfg.Namespace != "" {
		namespace = cfg.Namespace
		release, _ = chrt.CreateInstallParams(cfg.BuildId)
		releaseSelector = fmt.Sprintf("%s=%s", cfg.ReleaseLabel, release)
		cleanup = func() {
			helm.DeleteRelease(namespace, release)
		}
	} else {
		release, namespace = chrt.CreateInstallParams(cfg.BuildId)
		cleanup = func() {
			helm.DeleteRelease(namespace, release)
			kubectl.DeleteNamespace(namespace)
		}
	}
	return
}

// testRelease tests a release.
func testRelease(
	helm tool.Helm,
	kubectl tool.Kubectl,
	release, namespace, releaseSelector string,
	cleanupHelmTests bool,
) error {
	if err := kubectl.WaitForDeployments(namespace, releaseSelector); err != nil {
		return err
	}
	if err := helm.Test(namespace, release); err != nil {
		return err
	}
	return nil
}

// getChartPreviousVersion attemtps to retrieve the previous version
// of the given chart.
func getChartPreviousVersion(chrt *chart.Chart) (*chart.Chart, error) {
	// TODO: decide which sources do we consider when searching for a
	//       previous version's candidate
	return chrt, nil
}

// upgradeAndTestChart performs the installation of the given oldChrt,
// and attempts to perform an upgrade from that state.
func upgradeAndTestChart(
	cfg config.Configuration,
	oldChrt, chrt *chart.Chart,
	helm tool.Helm,
	kubectl tool.Kubectl,
) chart.TestResult {

	// result contains the test result; please notice that each values
	// file in the chart's 'ci' folder will be installed and tested
	// and the first failure makes the test fail.
	result := chart.TestResult{Chart: chrt}

	valuesFiles := oldChrt.ValuesFilePathsForCI()
	if len(valuesFiles) == 0 {
		valuesFiles = append(valuesFiles, "")
	}
	for _, valuesFile := range valuesFiles {
		if valuesFile != "" {
			if cfg.SkipMissingValues && !chrt.HasCIValuesFile(valuesFile) {
				// TODO: do not assume STDOUT here; instead a writer
				//       should be given to be written to.
				fmt.Printf("Upgrade testing for values file '%s' skipped because a corresponding values file was not found in %s/ci", valuesFile, chrt.Path())
				continue
			}
		}

		// Use anonymous function. Otherwise deferred calls would pile up
		// and be executed in reverse order after the loop.
		fun := func() error {
			namespace, release, releaseSelector, cleanup := generateInstallConfig(cfg, oldChrt, helm, kubectl)
			defer cleanup()

			// Install previous version of chart. If installation fails, ignore this release.
			if err := helm.InstallWithValues(oldChrt.Path(), valuesFile, namespace, release); err != nil {
				return fmt.Errorf("Upgrade testing for release '%s' skipped because of previous revision installation error: %w", release, err)
			}
			if err := testRelease(helm, kubectl, release, namespace, releaseSelector, true); err != nil {
				return fmt.Errorf("Upgrade testing for release '%s' skipped because of previous revision testing error", release)
			}

			if err := helm.Upgrade(oldChrt.Path(), namespace, release); err != nil {
				return err
			}

			return testRelease(helm, kubectl, release, namespace, releaseSelector, false)
		}

		if err := fun(); err != nil {
			result.Error = err
			break
		}
	}

	return result
}

// readObjectFromYamlFile unmarshals the given filename and returns an object with its contents.
func readObjectFromYamlFile(filename string) (map[string]interface{}, error) {
	objBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("reading values file: %w", err)
	}

	var obj map[string]interface{}
	err = yaml.Unmarshal(objBytes, &obj)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling values file contents: %w", err)
	}

	return obj, nil
}

// writeObjectToTempYamlFile writes the given obj into a temporary file.
//
// It is responsibility of the caller to discard the file when finished using it.
func writeObjectToTempYamlFile(obj map[string]interface{}) (string, error) {
	objBytes, err := yaml.Marshal(obj)
	if err != nil {
		return "", fmt.Errorf("marshalling values file new contents: %w", err)
	}

	tempDirName, err := ioutil.TempDir(os.TempDir(), "chart-testing-*")
	if err != nil {
		return "", fmt.Errorf("creating temporary directory: %w", err)
	}

	filename := path.Join(tempDirName, "values.yaml")

	err = ioutil.WriteFile(filename, objBytes, fs.ModeExclusive)
	if err != nil {
		return "", fmt.Errorf("writing values file new contents: %w", err)
	}

	return filename, nil
}

// applyExtraValues applies the extra values provided into the given filename (a YAML file) and materializes its
// contents in the file returned by the function.
func applyExtraValues(filename string, extraValues map[string]interface{}) (string, error) {

	obj, err := readObjectFromYamlFile(filename)
	if err != nil {
		return "", fmt.Errorf("reading values file: %w", err)
	}

	err = mergo.MergeWithOverwrite(obj, extraValues)
	if err != nil {
		return "", fmt.Errorf("merging extra values: %w", err)
	}

	newValuesFile, err := writeObjectToTempYamlFile(obj)
	if err != nil {
		return "", fmt.Errorf("writing object to temporary location: %w", err)
	}

	return newValuesFile, nil
}

// installAndTestChartRelease installs and tests a chart release.
func installAndTestChartRelease(
	cfg config.Configuration,
	chrt *chart.Chart,
	helm tool.Helm,
	kubectl tool.Kubectl,
	extraValues map[string]interface{},
) chart.TestResult {

	// valuesFiles contains all the configurations that should be
	// executed; in other words, it performs a test matrix between
	// values files and tests.
	valuesFiles := chrt.ValuesFilePathsForCI()

	// Test with defaults if no values files are specified.
	if len(valuesFiles) == 0 {
		valuesFiles = append(valuesFiles, "")
	}

	result := chart.TestResult{Chart: chrt}

	for _, valuesFile := range valuesFiles {

		newValuesFile, err := applyExtraValues(valuesFile, extraValues)
		if err != nil {
			result.Error = fmt.Errorf("applying extra values: %w", err)
		}
		defer func() {
			os.Remove(newValuesFile)
		}()

		// Use anonymous function. Otherwise deferred calls would pile up
		// and be executed in reverse order after the loop.
		fun := func() error {
			namespace, release, releaseSelector, cleanup := generateInstallConfig(cfg, chrt, helm, kubectl)
			defer cleanup()

			if err := helm.InstallWithValues(chrt.Path(), newValuesFile, namespace, release); err != nil {
				return err
			}
			return testRelease(helm, kubectl, release, namespace, releaseSelector, false)
		}

		if err := fun(); err != nil {
			// fail fast approach; could be changed to best effort.
			result.Error = err
			break
		}
	}

	return result
}
