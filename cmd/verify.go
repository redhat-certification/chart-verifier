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
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/getter"

	"github.com/pkg/errors"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/redhat-certification/chart-verifier/pkg/chartverifier"
)

var Version = "1.0.0"

func init() {
	allChecks = chartverifier.DefaultRegistry().AllChecks()
}

//goland:noinspection GoUnusedGlobalVariable
var (
	// allChecks contains all available checks to be executed by the program.
	allChecks []string
	// enabledChecksFlag are the checks that should be performed, after the command initialization has happened.
	enabledChecksFlag []string
	// disabledChecksFlag are the checks that should not be performed.
	disabledChecksFlag []string
	// outputFormatFlag contains the output format the user has specified: default, yaml or json.
	outputFormatFlag string
	// setOverridesFlag contains the overrides the user has specified through the --set flag.
	setOverridesFlag []string
)

func filterChecks(set []string, subset []string, setEnabled bool, subsetEnabled bool) ([]string, error) {
	selected := make([]string, 0)
	seen := map[string]bool{}
	for _, v := range set {
		seen[v] = setEnabled
	}
	for _, v := range subset {
		if _, ok := seen[v]; !ok {
			return nil, errors.Errorf("check %q is unknown", v)
		}
		seen[v] = subsetEnabled
	}
	for k, v := range seen {
		if v {
			selected = append(selected, k)
		}
	}
	return selected, nil
}

func buildChecks(all, enabled, disabled []string) ([]string, error) {
	switch {
	case len(enabled) > 0 && len(disabled) > 0:
		return nil, errors.New("--enable and --disable can't be used at the same time")
	case len(enabled) > 0:
		return filterChecks(all, enabled, false, true)
	case len(disabled) > 0:
		return filterChecks(all, disabled, true, false)
	default:
		return all, nil
	}
}

// settings comes from Helm, to extract the same configuration values Helm uses.
var settings = cli.New()

// NewVerifyCmd creates ...
func NewVerifyCmd(config *viper.Viper) *cobra.Command {

	// opts contains command line options extracted from the environment.
	opts := &values.Options{}

	cmd := &cobra.Command{
		Use:   "verify <chart-uri>",
		Args:  cobra.ExactArgs(1),
		Short: "Verifies a Helm chart by checking some of its characteristics",
		RunE: func(cmd *cobra.Command, args []string) error {
			checks, err := buildChecks(allChecks, enabledChecksFlag, disabledChecksFlag)
			if err != nil {
				return err
			}

			// vals is a resulting map considering all the options the user has given.
			vals, err := opts.MergeValues(getter.All(settings))
			if err != nil {
				return err
			}

			certifier, err := chartverifier.
				NewCertifierBuilder().
				SetValues(vals).
				SetChecks(checks).
				SetConfig(config).
				SetOverrides(setOverridesFlag).
				SetToolVersion(Version).
				Build()

			if err != nil {
				return err
			}

			result, err := certifier.Certify(args[0])
			if err != nil {
				return err
			}

			if outputFormatFlag == "json" {
				b, err := json.Marshal(result)
				if err != nil {
					return err
				}

				cmd.Println(string(b))

			} else if outputFormatFlag == "yaml" {
				b, err := yaml.Marshal(result)
				if err != nil {
					return err
				}

				cmd.Println(string(b))
			} else {
				cmd.Print(result)
			}

			reportErr := chartverifier.
				NewReportBuilder().
				SetCertificate(&result).
				SetChartUri(args[0]).
				Generate()

			if reportErr != nil {
				cmd.Println("Report failure :" + reportErr.Error())
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringSliceVarP(&opts.ValueFiles, "chart-values", "f", nil, "specify values in a YAML file or a URL (can specify multiple)")

	cmd.Flags().StringSliceVarP(&opts.Values, "chart-set", "S", nil, "set values for the chart (can specify multiple or separate values with commas: key1=val1,key2=val2)")

	cmd.Flags().StringSliceVarP(&opts.StringValues, "chart-set-string", "X", nil, "set STRING values for the chart (can specify multiple or separate values with commas: key1=val1,key2=val2)")

	cmd.Flags().StringSliceVarP(&opts.FileValues, "chart-set-file", "F", nil, "set values from respective files specified via the command line (can specify multiple or separate values with commas: key1=path1,key2=path2)")

	cmd.Flags().StringSliceVarP(&enabledChecksFlag, "enable", "e", nil, "only the informed checks will be enabled")

	cmd.Flags().StringSliceVarP(&disabledChecksFlag, "disable", "x", nil, "all checks will be enabled except the informed ones")

	cmd.Flags().StringVarP(&outputFormatFlag, "output", "o", "", "the output format: default, json or yaml")

	cmd.Flags().StringSliceVarP(&setOverridesFlag, "set", "s", []string{}, "overrides a configuration, e.g: dummy.ok=false")

	return cmd
}

func init() {
	rootCmd.AddCommand(NewVerifyCmd(viper.GetViper()))
}
