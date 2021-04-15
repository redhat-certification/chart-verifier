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

package chartverifier

import (
	"github.com/redhat-certification/chart-verifier/pkg/chartverifier/checks"
	helmchart "helm.sh/helm/v3/pkg/chart"
)

var CertificateApiVersion = "v1"
var CertificateKind = "verify-report"

type OutcomeType string

const (
	MandatoryCheckType    checks.CheckType = "Mandatory"
	OptionalCheckType     checks.CheckType = "Optional"
	ExperimentalCheckType checks.CheckType = "Experimental"

	FailOutcomeType OutcomeType = "FAIL"
	PassOutcomeType OutcomeType = "PASS"
)

type Certificate struct {
	Apiversion string              `json:"apiversion" yaml:"apiversion"`
	Kind       string              `json:"kind" yaml:"kind"`
	Metadata   CertificateMetadata `json:"metadata" yaml:"metadata"`
	Results    []*CheckReport      `json:"results" yaml:"results"`
}

type CertificateMetadata struct {
	ToolMetadata ToolMetadata        `json:"tool" yaml:"tool"`
	ChartData    *helmchart.Metadata `json:"chart" yaml:"chart"`
	Overrides    string              `json: "chart-overrides" yaml:"chart-overrides"`
}

type ToolMetadata struct {
	Version                    string `json:"verifier-version" yaml:"verifier-version"`
	ChartUri                   string `json:"chart-uri" yaml:"chart-uri"`
	Digest                     string `json:"digest" yaml:"digest"`
	LastCertifiedTime          string `json:"lastCertifiedTime" yaml:"lastCertifiedTime"`
	CertifiedOpenShiftVersions string `json:"certifiedOpenShiftVersions" yaml:"certifiedOpenShiftVersions"`
}

type CheckReport struct {
	Check   string           `json:"check" yaml:"check"`
	Type    checks.CheckType `json:"type" yaml:"type"`
	Outcome OutcomeType      `json:"outcome" yaml:"outcome"`
	Reason  string           `json:"reason" yaml:"reason"`
}

func newCertificate() Certificate {

	certificate := Certificate{Apiversion: CertificateApiVersion, Kind: CertificateKind}
	certificate.Metadata = CertificateMetadata{}
	certificate.Metadata.ToolMetadata = ToolMetadata{}

	return certificate
}

func (c *Certificate) AddCheck(checkName string, checkType checks.CheckType) *CheckReport {
	newCheck := CheckReport{}
	newCheck.Check = checkName
	newCheck.Type = checkType
	newCheck.Outcome = PassOutcomeType
	c.Results = append(c.Results, &newCheck)
	return &newCheck
}

func (cr *CheckReport) SetResult(outcome bool, reason string) {
	if outcome {
		cr.Outcome = PassOutcomeType
	} else {
		cr.Outcome = FailOutcomeType
	}
	cr.Reason = reason
}

func (c *Certificate) IsOk() bool {

	outcome := true
	for _, check := range c.Results {
		if check.Outcome == FailOutcomeType {
			outcome = false
			break
		}
	}
	return outcome
}
