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

	"github.com/pkg/errors"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/redhat-certification/chart-verifier/pkg/chartverifier"
)

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

func buildCertifier(checks []string) (chartverifier.Certifier, error) {
	return chartverifier.NewCertifierBuilder().
		SetChecks(checks).
		Build()
}

func NewCertifyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "certify <chart-uri>",
		Args:  cobra.ExactArgs(1),
		Short: "Certifies a Helm chart by checking some of its characteristics",
		RunE: func(cmd *cobra.Command, args []string) error {
			checks, err := buildChecks(allChecks, enabledChecksFlag, disabledChecksFlag)
			if err != nil {
				return err
			}

			certifier, err := buildCertifier(checks)
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

			return nil
		},
	}

	cmd.Flags().StringSliceVarP(&enabledChecksFlag, "enable", "e", nil, "only the informed checks will be enabled")

	cmd.Flags().StringSliceVarP(&disabledChecksFlag, "disable", "x", nil, "all checks will be enabled except the informed ones")

	cmd.Flags().StringVarP(&outputFormatFlag, "output", "o", "", "the output format: default, json or yaml")

	return cmd
}

// verifyCmd represents the lint command
var certifyCmd = NewCertifyCmd()

func init() {
	rootCmd.AddCommand(certifyCmd)
}
