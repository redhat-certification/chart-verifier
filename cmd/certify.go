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
	// chartUri contains the chart location as informed by the user; should accept anything that Helm understands as a Chart
	// URI.
	chartUri string
	// onlyChecks are the checks that should be performed, after the command initialization has happened.
	onlyChecks []string
	// exceptChecks are the checks that should not be performed.
	exceptChecks []string
	// outputFormat contains the output format the user has specified: default, yaml or json.
	outputFormat string
)

func buildChecks(allChecks, onlyChecks, _ []string) []string {
	if onlyChecks != nil {
		return onlyChecks
	}
	return allChecks
}

func buildCertifier(checks []string) (chartverifier.Certifier, error) {
	return chartverifier.NewCertifierBuilder().
		SetChecks(checks).
		Build()
}

func NewCertifyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "certify",
		Args:  cobra.NoArgs,
		Short: "Certifies a Helm chart by checking some of its characteristics",
		RunE: func(cmd *cobra.Command, args []string) error {

			checks := buildChecks(allChecks, onlyChecks, exceptChecks)

			certifier, err := buildCertifier(checks)
			if err != nil {
				return err
			}

			result, err := certifier.Certify(chartUri)
			if err != nil {
				return err
			}

			if outputFormat == "json" {
				b, err := json.Marshal(result)
				if err != nil {
					return err
				}

				cmd.Println(string(b))

			} else if outputFormat == "yaml" {
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

	cmd.Flags().StringVarP(&chartUri, "uri", "u", "", "uri of the Chart being certified")
	_ = cmd.MarkFlagRequired("uri")

	cmd.Flags().StringSliceVarP(&onlyChecks, "only", "o", nil, "only the informed checks will be performed")

	cmd.Flags().StringSliceVarP(&exceptChecks, "except", "e", nil, "all available checks except those informed will be performed")

	cmd.Flags().StringVarP(&outputFormat, "output", "f", "", "the output format: default, json or yaml")

	return cmd
}

// certifyCmd represents the lint command
var certifyCmd = NewCertifyCmd()

func init() {
	rootCmd.AddCommand(certifyCmd)
}
