/*
 * Copyright 2021 Red Hat
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/redhat-certification/chart-verifier/pkg/chartverifier/checks"
	"github.com/redhat-certification/chart-verifier/pkg/chartverifier/profiles"
	"github.com/redhat-certification/chart-verifier/pkg/chartverifier/utils"

	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/getter"

	"github.com/pkg/errors"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/redhat-certification/chart-verifier/pkg/chartverifier"
)

func init() {
	allChecks = chartverifier.DefaultRegistry().AllChecks()
}

//goland:noinspection GoUnusedGlobalVariable
var (
	// allChecks contains all available checks to be executed by the program.
	allChecks checks.DefaultRegistry
	// enabledChecksFlag are the checks that should be performed, after the command initialization has happened.
	enabledChecksFlag []string
	// disabledChecksFlag are the checks that should not be performed.
	disabledChecksFlag []string
	// outputFormatFlag contains the output format the user has specified: default, yaml or json.
	outputFormatFlag string
	// setOverridesFlag contains the overrides the user has specified through the --set flag.
	setOverridesFlag []string
	// openshiftVersionFlag set the value of `certifiedOpenShiftVersions` in the report
	openshiftVersionFlag string
	// write report to file
	reportToFile bool
	// write an error log file
	suppressErrorLog bool
)

func filterChecks(set profiles.FilteredRegistry, subset []string, setEnabled bool, subsetEnabled bool) (chartverifier.FilteredRegistry, error) {

	selected := make(chartverifier.FilteredRegistry, 0)
	seen := map[checks.CheckName]bool{}
	for k, _ := range set {
		seen[k] = setEnabled
	}
	for _, v := range subset {
		if _, ok := seen[checks.CheckName(v)]; !ok {
			return nil, errors.Errorf("check %q is unknown", v)
		}
		seen[checks.CheckName(v)] = subsetEnabled
	}
	for k, v := range seen {
		if v {
			selected[k] = set[k]
		}
	}
	return selected, nil
}

func buildChecks(all checks.DefaultRegistry, config *viper.Viper, enabled, disabled []string) (chartverifier.FilteredRegistry, error) {
	profileChecks := profiles.New(config).FilterChecks(all)
	switch {
	case len(enabled) > 0 && len(disabled) > 0:
		return nil, errors.New("--enable and --disable can't be used at the same time")
	case len(enabled) > 0:
		return filterChecks(profileChecks, enabled, false, true)
	case len(disabled) > 0:
		return filterChecks(profileChecks, disabled, true, false)
	default:
		return chartverifier.FilteredRegistry(profileChecks), nil
	}
}

// settings comes from Helm, to extract the same configuration values Helm uses.
var settings = cli.New()

type verifyOptions struct {
	ValueFiles []string
	Values     []string
}

// NewVerifyCmd creates ...
func NewVerifyCmd(config *viper.Viper) *cobra.Command {

	// opts contains command line options extracted from the environment.
	opts := &values.Options{}

	// verifyOpts contains this specific command options.
	verifyOpts := &verifyOptions{}

	cmd := &cobra.Command{
		Use:   "verify <chart-uri>",
		Args:  cobra.ExactArgs(1),
		Short: "Verifies a Helm chart by checking some of its characteristics",
		RunE: func(cmd *cobra.Command, args []string) error {

			reportName := ""
			if reportToFile {
				if outputFormatFlag == "json" {
					reportName = "report.json"
				} else {
					reportName = "report.yaml"
				}
			}
			utils.InitLog(cmd, reportName, suppressErrorLog)

			utils.LogInfo(fmt.Sprintf("Chart Verifer %s.", Version))
			utils.LogInfo(fmt.Sprintf("Verify : %s", args[0]))

			// vals is a resulting map considering all the options the user has given.
			vals, err := opts.MergeValues(getter.All(settings))
			if err != nil {
				utils.LogError(err.Error())
				return err
			}

			verifierBuilder := chartverifier.NewVerifierBuilder().
				SetValues(vals).
				SetConfig(config).
				SetOverrides(verifyOpts.Values)

			checks, err := buildChecks(allChecks, verifierBuilder.GetConfig(), enabledChecksFlag, disabledChecksFlag)
			if err != nil {
				utils.LogError(err.Error())
				return err
			}

			verifier, err := verifierBuilder.
				SetChecks(checks).
				SetToolVersion(Version).
				SetOpenShiftVersion(openshiftVersionFlag).
				Build()

			if err != nil {
				return err
			}

			result, err := verifier.Verify(args[0])
			if err != nil {
				utils.LogError(err.Error())
				return err
			}

			outputContent := ""
			if outputFormatFlag == "json" {
				b, err := json.Marshal(result)
				if err != nil {
					utils.LogError(err.Error())
					return err
				}
				outputContent = string(b)
			} else {
				b, err := yaml.Marshal(result)
				if err != nil {
					return err
				}
				outputContent = string(b)
			}
			utils.WriteStdOut(outputContent)

			utils.WriteLogs(outputFormatFlag)

			return nil
		},
	}

	settings.AddFlags(cmd.Flags())

	cmd.Flags().StringSliceVarP(&opts.ValueFiles, "chart-values", "F", nil, "specify values in a YAML file or a URL (can specify multiple)")

	cmd.Flags().StringSliceVarP(&opts.Values, "chart-set", "S", nil, "set values for the chart (can specify multiple or separate values with commas: key1=val1,key2=val2)")

	cmd.Flags().StringSliceVarP(&opts.StringValues, "chart-set-string", "X", nil, "set STRING values for the chart (can specify multiple or separate values with commas: key1=val1,key2=val2)")

	cmd.Flags().StringSliceVarP(&opts.FileValues, "chart-set-file", "G", nil, "set values from respective files specified via the command line (can specify multiple or separate values with commas: key1=path1,key2=path2)")

	cmd.Flags().StringSliceVarP(&enabledChecksFlag, "enable", "e", nil, "only the informed checks will be enabled")

	cmd.Flags().StringSliceVarP(&disabledChecksFlag, "disable", "x", nil, "all checks will be enabled except the informed ones")

	cmd.Flags().StringVarP(&outputFormatFlag, "output", "o", "", "the output format: default, json or yaml")

	cmd.Flags().StringSliceVarP(&verifyOpts.Values, "set", "s", []string{}, "overrides a configuration, e.g: dummy.ok=false")

	cmd.Flags().StringSliceVarP(&verifyOpts.ValueFiles, "set-values", "f", nil, "specify application and check configuration values in a YAML file or a URL (can specify multiple)")
	cmd.Flags().StringVarP(&openshiftVersionFlag, "openshift-version", "V", "", "version of OpenShift used in the cluster")
	cmd.Flags().BoolVarP(&reportToFile, "write-to-file", "w", false, "write report to report.yaml (default: stdout)")
	cmd.Flags().BoolVarP(&suppressErrorLog, "suppress-error-log", "E", false, "suppress the error log (default: written to ./chart-verifier.log)")
	return cmd
}

func init() {
	rootCmd.AddCommand(NewVerifyCmd(viper.GetViper()))
}
