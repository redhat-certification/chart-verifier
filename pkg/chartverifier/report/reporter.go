package report

import (
	"fmt"
	"github.com/redhat-certification/chart-verifier/pkg/chartverifier"
	"github.com/redhat-certification/chart-verifier/pkg/chartverifier/checks"
	"github.com/redhat-certification/chart-verifier/pkg/chartverifier/profiles"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	AnnotationsPrefixConfigName string = "annotations.prefix"

	DefaultAnnotationsPrefix string = "charts.openshift.io"

	DigestsAnnotationName                string = "digest"
	LastCertifiedTimestampAnnotationName string = "lastCertifiedTimestamp"
	CertifiedOCPVersionAnnotationName    string = "certifiedOpenShiftVersions"
)

func All(opts *ReportOptions) (OutputReport, error) {

	outputReport := OutputReport{}

	subReport, err := Annotations(opts)
	if err == nil {
		outputReport.AnnotationsReport = subReport.AnnotationsReport
	}

	subReport, err = Digests(opts)
	if err == nil {
		outputReport.DigestsReport = subReport.DigestsReport
	}

	subReport, err = Results(opts)
	if err == nil {
		outputReport.ResultsReport = subReport.ResultsReport
	}

	subReport, err = Metadata(opts)
	if err == nil {
		outputReport.MetadataReport = subReport.MetadataReport
	}

	return outputReport, nil

}

func Annotations(opts *ReportOptions) (OutputReport, error) {

	var outputReport OutputReport
	report, err := readReport(opts.URI)
	if err != nil {
		return outputReport, err
	}

	outputReport = OutputReport{}
	anotationsPrefix := DefaultAnnotationsPrefix

	if opts.ViperConfig != nil {
		configAnnotationsPrefix := opts.ViperConfig.GetString(AnnotationsPrefixConfigName)
		if len(configAnnotationsPrefix) > 0 {
			anotationsPrefix = configAnnotationsPrefix
		}
	}

	name := fmt.Sprintf("%s/%s", anotationsPrefix, DigestsAnnotationName)
	value := report.Metadata.ToolMetadata.Digests.Chart
	if len(value) > 0 {
		annotation := Annotation{}
		annotation.Name = name
		annotation.Value = value
		outputReport.AnnotationsReport = append(outputReport.AnnotationsReport, annotation)
	}

	name = fmt.Sprintf("%s/%s", anotationsPrefix, LastCertifiedTimestampAnnotationName)
	value = report.Metadata.ToolMetadata.LastCertifiedTimestamp
	if len(value) > 0 {
		annotation := Annotation{}
		annotation.Name = name
		annotation.Value = value
		outputReport.AnnotationsReport = append(outputReport.AnnotationsReport, annotation)
	}

	name = fmt.Sprintf("%s/%s", anotationsPrefix, CertifiedOCPVersionAnnotationName)
	value = report.Metadata.ToolMetadata.CertifiedOpenShiftVersions
	if len(value) > 0 {
		annotation := Annotation{}
		annotation.Name = name
		annotation.Value = value
		outputReport.AnnotationsReport = append(outputReport.AnnotationsReport, annotation)
	}

	return outputReport, nil
}

func Digests(opts *ReportOptions) (OutputReport, error) {
	var outputReport OutputReport

	report, err := readReport(opts.URI)
	if err != nil {
		return outputReport, err
	}

	outputReport = OutputReport{}
	outputReport.DigestsReport = &DigestReport{}

	outputReport.DigestsReport.ChartDigest = report.Metadata.ToolMetadata.Digests.Chart
	outputReport.DigestsReport.PackageDigest = report.Metadata.ToolMetadata.Digests.Package

	return outputReport, nil

}

func Metadata(opts *ReportOptions) (OutputReport, error) {

	var outputReport OutputReport

	report, err := readReport(opts.URI)
	if err != nil {
		return outputReport, err
	}

	outputReport = OutputReport{}
	outputReport.MetadataReport = &MetadataReport{}

	outputReport.MetadataReport.ProfileVendorType = profiles.VendorType(report.Metadata.ToolMetadata.Profile.VendorType)
	outputReport.MetadataReport.ProfileVersion = report.Metadata.ToolMetadata.Profile.Version
	return outputReport, nil

}

func Results(opts *ReportOptions) (OutputReport, error) {

	var outputReport OutputReport
	var messages []string

	report, err := readReport(opts.URI)
	if err != nil {
		return outputReport, err
	}

	profileVendorType := report.Metadata.ToolMetadata.Profile.VendorType
	profileVersion := report.Metadata.ToolMetadata.Profile.VendorType

	if opts.ViperConfig != nil {

		configVendorType := profiles.VendorType(opts.ViperConfig.GetString(profiles.VendorTypeConfigName))
		if len(configVendorType) > 0 {
			profileVendorType = string(configVendorType)
		}

		configProfileVersion := opts.ViperConfig.GetString(profiles.VersionConfigName)
		if len(configProfileVersion) > 0 {
			profileVersion = configProfileVersion
		}
	}

	config := viper.New()
	config.Set(profiles.VendorTypeConfigName, profileVendorType)
	config.Set(profiles.VersionConfigName, profileVersion)
	profile := profiles.New(config)

	passed := 0
	failed := 0

	for _, profileCheck := range profile.Checks {
		if profileCheck.Type == checks.MandatoryCheckType {
			found := false
			for _, reportCheck := range report.Results {
				if strings.Compare(profileCheck.Name, string(reportCheck.Check)) == 0 {
					found = true
					if reportCheck.Outcome == chartverifier.PassOutcomeType {
						passed++
					} else {
						failed++
						// Change multiple line reasons to a single line
						reason := strings.ReplaceAll(strings.TrimRight(reportCheck.Reason, "\n"), "\n", ", ")
						messages = append(messages, reason)
					}
					break
				}
			}
			if !found {
				failed++
				messages = append(messages, fmt.Sprintf("Missing mandatory check : %s", profileCheck.Name))
			}
		}
	}

	outputReport = OutputReport{}
	outputReport.ResultsReport = &ResultsReport{}

	outputReport.ResultsReport.Passed = fmt.Sprintf("%d", passed)
	outputReport.ResultsReport.Failed = fmt.Sprintf("%d", failed)
	outputReport.ResultsReport.Messages = messages

	return outputReport, nil

}

var loadedReport *reportInfo

type reportInfo struct {
	uri    string
	report *chartverifier.Report
}

func readReport(path string) (*chartverifier.Report, error) {

	if loadedReport == nil || loadedReport.report == nil {
		loadedReport = &reportInfo{}
		loadedReport.uri = path
	} else if strings.Compare(path, loadedReport.uri) != 0 {
		loadedReport.uri = path
		loadedReport.report = nil
	} else {
		return loadedReport.report, nil
	}

	reportPath, err := filepath.Abs(path)

	// Open the yaml file which defines the tests to run
	reportYaml, err := os.Open(reportPath)
	if err != nil {
		return nil, err
	}

	reportBytes, err := ioutil.ReadAll(reportYaml)
	if err != nil {
		return nil, err
	}

	loadedReport.report = &chartverifier.Report{}
	err = yaml.Unmarshal(reportBytes, loadedReport.report)
	if err != nil {
		return nil, err
	}

	return loadedReport.report, nil

}
