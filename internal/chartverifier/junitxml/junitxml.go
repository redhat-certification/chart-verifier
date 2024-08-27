package junitxml

import (
	"encoding/xml"
	"fmt"
	"strconv"

	"github.com/redhat-certification/chart-verifier/pkg/chartverifier/report"
)

type JUnitTestSuites struct {
	XMLName xml.Name         `xml:"testsuites"`
	Suites  []JUnitTestSuite `xml:"testsuite"`
}

type JUnitTestSuite struct {
	XMLName      xml.Name        `xml:"testsuite"`
	Tests        int             `xml:"tests,attr"`
	Failures     int             `xml:"failures,attr"`
	Skipped      int             `xml:"skipped,attr"`
	Unknown      int             `xml:"unknown,attr"`
	ReportDigest string          `xml:"reportDigest,attr"`
	Name         string          `xml:"name,attr"`
	Properties   []JUnitProperty `xml:"properties>property,omitempty"`
	TestCases    []JUnitTestCase `xml:"testcase"`
}

type JUnitTestCase struct {
	XMLName     xml.Name          `xml:"testcase"`
	Classname   string            `xml:"classname,attr"`
	Name        string            `xml:"name,attr"`
	SkipMessage *JUnitSkipMessage `xml:"skipped,omitempty"`
	Failure     *JUnitMessage     `xml:"failure,omitempty"`
	Warning     *JUnitMessage     `xml:"warning,omitempty"`
	SystemOut   string            `xml:"system-out,omitempty"`
	Message     string            `xml:",chardata"`
}

type JUnitSkipMessage struct {
	Message string `xml:"message,attr"`
}

type JUnitProperty struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type JUnitMessage struct {
	Message  string `xml:"message,attr"`
	Type     string `xml:"type,attr"`
	Contents string `xml:",chardata"`
}

func Format(r report.Report) ([]byte, error) {
	results := r.Results
	checksByOutcome := map[string][]report.CheckReport{}

	for i, result := range results {
		checksByOutcome[result.Outcome] = append(checksByOutcome[result.Outcome], *results[i])
	}

	digest, err := r.GetReportDigest()
	if err != nil {
		// Prefer to continue even if digest calculation fails for some reason.
		digest = "unknown"
	}

	testsuite := JUnitTestSuite{
		Tests:        len(results),
		Failures:     len(checksByOutcome[report.FailOutcomeType]),
		Skipped:      len(checksByOutcome[report.SkippedOutcomeType]),
		Unknown:      len(checksByOutcome[report.UnknownOutcomeType]),
		ReportDigest: digest,
		Name:         "Red Hat Helm Chart Certification",
		Properties: []JUnitProperty{
			{Name: "profileType", Value: r.Metadata.ToolMetadata.Profile.VendorType},
			{Name: "profileVersion", Value: r.Metadata.ToolMetadata.Profile.Version},
			{Name: "webCatalogOnly", Value: strconv.FormatBool(r.Metadata.ToolMetadata.ProviderDelivery || r.Metadata.ToolMetadata.WebCatalogOnly)},
			{Name: "verifierVersion", Value: r.Metadata.ToolMetadata.Version},
		},
		TestCases: []JUnitTestCase{},
	}

	for _, tc := range checksByOutcome[report.PassOutcomeType] {
		c := JUnitTestCase{
			Classname: r.Metadata.ToolMetadata.ChartUri,
			Name:      string(tc.Check),
			Failure:   nil,
			Message:   tc.Reason,
		}
		testsuite.TestCases = append(testsuite.TestCases, c)
	}

	for _, tc := range checksByOutcome[report.FailOutcomeType] {
		c := JUnitTestCase{
			Classname: r.Metadata.ToolMetadata.ChartUri,
			Name:      string(tc.Check),
			Failure: &JUnitMessage{
				Message:  "Failed",
				Type:     string(tc.Type),
				Contents: tc.Reason,
			},
			Message: tc.Reason,
		}
		testsuite.TestCases = append(testsuite.TestCases, c)
	}

	for _, tc := range checksByOutcome[report.UnknownOutcomeType] {
		c := JUnitTestCase{
			Classname: r.Metadata.ToolMetadata.ChartUri,
			Name:      string(tc.Check),
			Failure: &JUnitMessage{
				Message:  "Unknown",
				Type:     string(tc.Type),
				Contents: tc.Reason,
			},
			Message: tc.Reason,
		}
		testsuite.TestCases = append(testsuite.TestCases, c)
	}

	for _, tc := range checksByOutcome[report.SkippedOutcomeType] {
		c := JUnitTestCase{
			Classname: r.Metadata.ToolMetadata.ChartUri,
			Name:      string(tc.Check),
			Failure:   nil,
			Message:   tc.Reason,
			SkipMessage: &JUnitSkipMessage{
				Message: tc.Reason,
			},
		}
		testsuite.TestCases = append(testsuite.TestCases, c)
	}

	suites := JUnitTestSuites{
		Suites: []JUnitTestSuite{testsuite},
	}

	bytes, err := xml.MarshalIndent(suites, "", "\t")
	if err != nil {
		o := fmt.Errorf("error formatting results with formatter %s: %v",
			"junitxml",
			err,
		)

		return nil, o
	}

	return bytes, nil
}
