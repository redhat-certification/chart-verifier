package checks

import (
	"fmt"
	"strings"

	"github.com/helm/chart-testing/pkg/chart"
	"github.com/helm/chart-testing/pkg/config"
	"github.com/helm/chart-testing/pkg/exec"
	"github.com/helm/chart-testing/pkg/tool"
	"github.com/helm/chart-testing/pkg/util"
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

	if len(cfg.ReleaseLabel) == 0 {
		cfg.ReleaseLabel = "app.kubernetes.io/instance"
	}

	if len(cfg.Namespace) == 0 {
		cfg.Namespace = opts.HelmEnvSettings.Namespace()
	}

	return cfg
}

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
		result := upgradeAndTestChartFromPreviousRelease(cfg, chrt, helm, kubectl)
		if result.Error != nil {
			return NewResult(false, result.Error.Error()), nil
		}
	} else {
		result := installAndTestChartRelease(cfg, chrt, helm, kubectl)
		if result.Error != nil {
			return NewResult(false, result.Error.Error()), nil
		}
	}

	return NewResult(true, ""), nil
}

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
			helm.DeleteRelease(release)
		}
	} else {
		release, namespace = chrt.CreateInstallParams(cfg.BuildId)
		cleanup = func() {
			helm.DeleteRelease(release)
			kubectl.DeleteNamespace(namespace)
		}
	}
	return
}

// testRelease tests a chart release.
func testRelease(
	helm tool.Helm,
	kubectl tool.Kubectl,
	release, namespace, releaseSelector string,
	cleanupHelmTests bool,
) error {
	if err := kubectl.WaitForDeployments(namespace, releaseSelector); err != nil {
		return err
	}
	if err := helm.Test(release, cleanupHelmTests); err != nil {
		return err
	}
	return nil
}

// getChartPreviousVersion attemtps to retrieve the previous version
// of the given chart.
func getChartPreviousVersion(chrt *chart.Chart) (*chart.Chart, error) {
	return chrt, nil

}

func upgradeAndTestChartFromPreviousRelease(
	cfg config.Configuration,
	chrt *chart.Chart,
	helm tool.Helm,
	kubectl tool.Kubectl,
) chart.TestResult {
	result := chart.TestResult{Chart: chrt}

	oldChrt, err := getChartPreviousVersion(chrt)
	if err != nil {
		result.Error = fmt.Errorf("skipping upgrade test of '%s' because no previous chart is available", chrt.Yaml().Name)
		return result
	}

	breakingChangeAllowed, err := util.BreakingChangeAllowed(oldChrt.Yaml().Version, chrt.Yaml().Version)
	if !breakingChangeAllowed {
		result.Error = fmt.Errorf("Skipping upgrade test of '%s' because breaking changes are not allowed for chart", chrt)
		return result
	} else if err != nil {
		result.Error = err
		return result
	}

	result.Error = upgradeChart(cfg, oldChrt, chrt, helm, kubectl)
	return result
}

func upgradeChart(
	cfg config.Configuration,
	oldChrt, chrt *chart.Chart,
	helm tool.Helm,
	kubectl tool.Kubectl,
) error {
	return nil
}

// installAndTestChartRelease installs and tests a chart release.
func installAndTestChartRelease(
	cfg config.Configuration,
	chrt *chart.Chart,
	helm tool.Helm,
	kubectl tool.Kubectl,
) chart.TestResult {
	fmt.Printf("Installing chart '%s'...\n", chrt)
	valuesFiles := chrt.ValuesFilePathsForCI()

	// Test with defaults if no values files are specified.
	if len(valuesFiles) == 0 {
		valuesFiles = append(valuesFiles, "")
	}

	result := chart.TestResult{Chart: chrt}

	for _, valuesFile := range valuesFiles {
		if valuesFile != "" {
			fmt.Printf("\nInstalling chart with values file '%s'...\n\n", valuesFile)
		}

		// Use anonymous function. Otherwise deferred calls would pile up
		// and be executed in reverse order after the loop.
		fun := func() error {
			namespace, release, releaseSelector, cleanup := generateInstallConfig(cfg, chrt, helm, kubectl)
			defer cleanup()

			if err := helm.InstallWithValues(chrt.Path(), valuesFile, namespace, release); err != nil {
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
