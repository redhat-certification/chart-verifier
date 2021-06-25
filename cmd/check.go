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
	"io/ioutil"
	"regexp"

	"github.com/redhat-certification/chart-verifier/pkg/chartverifier"
	"github.com/redhat-certification/chart-verifier/pkg/chartverifier/profiles"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

const (
	defaultProfile   = "default"
	partnerProfile   = "partner"
	redhatProfile    = "redhat"
	communityProfile = "community"
)

type checkResult struct {
	Profile     string                       `json:"profile,omitempty" yaml:"profile,omitempty"`
	Passed      int                          `json:"passed" yaml:"passed"`
	Failed      int                          `json:"failed" yaml:"failed"`
	Unknown     int                          `json:"unknown" yaml:"unknown"`
	Message     []*chartverifier.CheckReport `json:"message" yaml:"message"`
	OtherResult *checkResult                 `json:"other-result,omitempty" yaml:"other-result,omitempty"`
}

type checkOptions struct {
	Profile string
}

func init() {
	rootCmd.AddCommand(NewCheckCmd(viper.GetViper()))
}

func checkReport(path string, profile string) (checkResult, error) {
	result := checkResult{}
	reportBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return result, fmt.Errorf("reading report file: %w", err)
	}

	var report chartverifier.Report
	err = yaml.Unmarshal(reportBytes, &report)
	if err != nil {
		return result, fmt.Errorf("unmarshalling report file contents: %w", err)
	}

	// TODO: use different profiles based on command line argument
	var p *profiles.Profile
	switch profile {
	case partnerProfile:
		p = profiles.Get()
		result.Profile = partnerProfile
	case redhatProfile:
		p = profiles.Get()
		result.Profile = redhatProfile
	case communityProfile:
		p = profiles.Get()
		result.Profile = communityProfile
	default:
		p = profiles.Get()
		result.Profile = defaultProfile
	}

	profileChecks := p.Checks
	profileChecksSet := make(map[string]bool)
	for _, c := range profileChecks {
		splitter := regexp.MustCompile(`/`)
		splitCheck := splitter.Split(c.Name, -1)
		profileChecksSet[splitCheck[1]] = true
	}

	otherResult := checkResult{}

	for _, cr := range report.Results {
		if _, ok := profileChecksSet[string(cr.Check)]; ok {
			switch cr.Outcome {
			case chartverifier.FailOutcomeType:
				result.Failed++
				result.Message = append(result.Message, cr)
			case chartverifier.PassOutcomeType:
				result.Passed++
			case chartverifier.UnknownOutcomeType:
				result.Unknown++
				result.Message = append(result.Message, cr)
			default:
				return result, fmt.Errorf("checking report results: incorrect outcome type '%v' for '%v'", cr.Outcome, cr.Check)
			}
		} else {
			switch cr.Outcome {
			case chartverifier.FailOutcomeType:
				otherResult.Failed++
				otherResult.Message = append(otherResult.Message, cr)
			case chartverifier.PassOutcomeType:
				otherResult.Passed++
			case chartverifier.UnknownOutcomeType:
				otherResult.Unknown++
				otherResult.Message = append(otherResult.Message, cr)
			default:
				return result, fmt.Errorf("checking report results: incorrect outcome type '%v' for '%v'", cr.Outcome, cr.Check)
			}
		}
	}

	if otherResult.Passed+otherResult.Failed+otherResult.Unknown > 0 {
		result.OtherResult = &otherResult
	}

	return result, nil
}

// NewCheckCmd creates a command that sanity checks report.
func NewCheckCmd(config *viper.Viper) *cobra.Command {

	// verifyOpts contains this specific command options.
	checkOpts := &checkOptions{}

	cmd := &cobra.Command{
		Use:   "check <report-uri>",
		Args:  cobra.ExactArgs(1),
		Short: "Checks the result of a Helm chart report",
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := checkReport(args[0], checkOpts.Profile)
			if err != nil {
				return err
			}

			if outputFormatFlag == "json" {
				b, err := json.Marshal(result)
				if err != nil {
					return err
				}

				cmd.Println(string(b))

			} else {
				b, err := yaml.Marshal(result)
				if err != nil {
					return err
				}

				cmd.Println(string(b))
			}
			return nil
		},
	}

	settings.AddFlags(cmd.Flags())

	cmd.Flags().StringVarP(&outputFormatFlag, "output", "o", "", "the output format: default, json or yaml")

	cmd.Flags().StringVarP(&checkOpts.Profile, "profile", "p", defaultProfile, "check according to specific profile (partner, redhat, community)")

	return cmd
}
