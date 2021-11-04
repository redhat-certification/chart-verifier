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
	"fmt"
	"github.com/redhat-certification/chart-verifier/pkg/chartverifier/checks"
	helmchart "helm.sh/helm/v3/pkg/chart"
)

var ReportApiVersion = "v1"
var ReportKind = "verify-report"

type OutcomeType string

const (
	FailOutcomeType    OutcomeType = "FAIL"
	PassOutcomeType    OutcomeType = "PASS"
	UnknownOutcomeType OutcomeType = "UNKNOWN"
)

type Report struct {
	Apiversion string         `json:"apiversion" yaml:"apiversion"`
	Kind       string         `json:"kind" yaml:"kind"`
	Metadata   ReportMetadata `json:"metadata" yaml:"metadata"`
	Results    []*CheckReport `json:"results" yaml:"results"`
}

type ReportMetadata struct {
	ToolMetadata ToolMetadata        `json:"tool" yaml:"tool"`
	ChartData    *helmchart.Metadata `json:"chart" yaml:"chart"`
	Overrides    string              `json:"chart-overrides" yaml:"chart-overrides"`
}

type ToolMetadata struct {
	Version                    string  `json:"verifier-version" yaml:"verifier-version"`
	Profile                    Profile `json:"profile" yaml:"profile"`
	ChartUri                   string  `json:"chart-uri" yaml:"chart-uri"`
	Digests                    Digests `json:"digests" yaml:"digests"`
	LastCertifiedTimestamp     string  `json:"lastCertifiedTimestamp,omitempty" yaml:"lastCertifiedTimestamp,omitempty"`
	CertifiedOpenShiftVersions string  `json:"certifiedOpenShiftVersions,omitempty" yaml:"certifiedOpenShiftVersions,omitempty"`
}

type Digests struct {
	Chart   string `json:"chart" yaml:"chart"`
	Package string `json:"package,omitempty" yaml:"package,omitempty"`
}

type Profile struct {
	VendorType string `json:"vendorType" yaml:"VendorType"`
	Version    string `json:"version" yaml:"version"`
}

type CheckReport struct {
	Check   checks.CheckName `json:"check" yaml:"check"`
	Type    checks.CheckType `json:"type" yaml:"type"`
	Outcome OutcomeType      `json:"outcome" yaml:"outcome"`
	Reason  string           `json:"reason" yaml:"reason"`
}

func newReport() Report {

	report := Report{Apiversion: ReportApiVersion, Kind: ReportKind}
	report.Metadata = ReportMetadata{}
	report.Metadata.ToolMetadata = ToolMetadata{}

	return report
}

func (c *Report) AddCheck(check checks.Check) *CheckReport {
	newCheck := CheckReport{}
	newCheck.Check = checks.CheckName(fmt.Sprintf("%s/%s", check.CheckId.Version, check.CheckId.Name))
	newCheck.Type = check.Type
	newCheck.Outcome = UnknownOutcomeType
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
