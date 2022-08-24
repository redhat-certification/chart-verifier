package reportsummary

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	apireport "github.com/redhat-certification/chart-verifier/pkg/chartverifier/report"
)

func TestAddStringReport(t *testing.T) {

	chartUri := "https://github.com/redhat-certification/chart-verifier/blob/main/tests/charts/psql-service/0.1.9/psql-service-0.1.9.tgz?raw=true"
	yamlFileReport := "test-reports/report.yaml"
	jsonFileReport := "test-reports/report.json"

	yamlFileReportBytes, yamlReportErr := loadChartFromAbsPath(yamlFileReport)
	require.NoError(t, yamlReportErr)
	jsonFileReportBytes, jsonReportErr := loadChartFromAbsPath(jsonFileReport)
	require.NoError(t, jsonReportErr)

	report, loadErr := apireport.NewReport().SetContent(string(yamlFileReportBytes)).Load()
	require.NoError(t, loadErr)
	reportSummary := NewReportSummary().SetReport(report)
	checkReportSummaries(reportSummary, chartUri, t)

	report, loadErr = apireport.NewReport().SetContent(string(jsonFileReportBytes)).Load()
	require.NoError(t, loadErr)
	reportSummary = NewReportSummary().SetReport(report)
	checkReportSummaries(reportSummary, chartUri, t)

}

func TestBadDigestReport(t *testing.T) {

	yamlFileReport := "test-reports/baddigest/report.yaml"
	jsonFileReport := "test-reports/baddigest/report.json"

	yamlFileReportBytes, yamlReportErr := loadChartFromAbsPath(yamlFileReport)
	require.NoError(t, yamlReportErr)
	jsonFileReportBytes, jsonReportErr := loadChartFromAbsPath(jsonFileReport)
	require.NoError(t, jsonReportErr)

	_, loadErr := apireport.NewReport().SetContent(string(yamlFileReportBytes)).Load()
	require.Error(t, loadErr)
	require.Contains(t, fmt.Sprintf("%v", loadErr), "Digest in report did not match report content")

	_, loadErr = apireport.NewReport().SetContent(string(jsonFileReportBytes)).Load()
	require.Error(t, loadErr)
	require.Contains(t, fmt.Sprintf("%v", loadErr), "Digest in report did not match report content")

}

func TestMissingDigestReport(t *testing.T) {

	yamlFileReport := "test-reports/missingdigest/report.yaml"
	jsonFileReport := "test-reports/missingdigest/report.json"

	yamlFileReportBytes, yamlReportErr := loadChartFromAbsPath(yamlFileReport)
	require.NoError(t, yamlReportErr)
	jsonFileReportBytes, jsonReportErr := loadChartFromAbsPath(jsonFileReport)
	require.NoError(t, jsonReportErr)

	_, loadErr := apireport.NewReport().SetContent(string(yamlFileReportBytes)).Load()
	require.Error(t, loadErr)
	require.Contains(t, fmt.Sprintf("%v", loadErr), "Report does not contain expected report digest.")

	_, loadErr = apireport.NewReport().SetContent(string(jsonFileReportBytes)).Load()
	require.Error(t, loadErr)
	require.Contains(t, fmt.Sprintf("%v", loadErr), "Report does not contain expected report digest.")

}

func TestPreDigestReport(t *testing.T) {

	chartUri := "https://github.com/redhat-certification/chart-verifier/blob/main/tests/charts/psql-service/0.1.9/psql-service-0.1.9.tgz?raw=true"
	yamlFileReport := "test-reports/predigest/report.yaml"
	jsonFileReport := "test-reports/predigest/report.json"

	yamlFileReportBytes, yamlReportErr := loadChartFromAbsPath(yamlFileReport)
	require.NoError(t, yamlReportErr)
	jsonFileReportBytes, jsonReportErr := loadChartFromAbsPath(jsonFileReport)
	require.NoError(t, jsonReportErr)

	report, loadErr := apireport.NewReport().SetContent(string(yamlFileReportBytes)).Load()
	require.NoError(t, loadErr)
	reportSummary := NewReportSummary().SetReport(report)
	checkReportSummaries(reportSummary, chartUri, t)

	report, loadErr = apireport.NewReport().SetContent(string(jsonFileReportBytes)).Load()
	require.NoError(t, loadErr)
	reportSummary = NewReportSummary().SetReport(report)
	checkReportSummaries(reportSummary, chartUri, t)

}

func TestAddUrlReport(t *testing.T) {

	yamlUrlReport := "https://github.com/redhat-certification/chart-verifier/blob/main/cmd/test/report.yaml?raw=true"
	url, loadUrlErr := url.Parse(yamlUrlReport)
	require.NoError(t, loadUrlErr)

	chartUri := "internal/chartverifier/checks/chart-0.1.0-v3.valid.tgz"
	report, loadErr := apireport.NewReport().SetURL(url).Load()
	require.NoError(t, loadErr)
	reportSummary := NewReportSummary().SetReport(report)

	checkReportSummaries(reportSummary, chartUri, t)

}
func checkReportSummaries(summary APIReportSummary, chartUri string, t *testing.T) {
	checkReportSummariesFormat(YamlReport, summary, chartUri, t)
	checkReportSummariesFormat(JsonReport, summary, chartUri, t)
}

func checkReportSummariesFormat(format SummaryFormat, summary APIReportSummary, chartUri string, t *testing.T) {

	reportSummary, reportSummaryErr := summary.GetContent(AllSummary, format)
	require.NoError(t, reportSummaryErr)
	checkSummaryAll(format, reportSummary, chartUri, t)

	reportSummary, reportSummaryErr = summary.GetContent(DigestsSummary, format)
	require.NoError(t, reportSummaryErr)
	checkSummaryDigests(false, format, reportSummary, chartUri, t)

	reportSummary, reportSummaryErr = summary.GetContent(MetadataSummary, format)
	require.NoError(t, reportSummaryErr)
	checkSummaryMetadata(false, format, reportSummary, chartUri, t)

	reportSummary, reportSummaryErr = summary.GetContent(AnnotationsSummary, format)
	require.NoError(t, reportSummaryErr)
	checkSummaryAnnotations(false, format, reportSummary, chartUri, t)

	reportSummary, reportSummaryErr = summary.GetContent(ResultsSummary, format)
	require.NoError(t, reportSummaryErr)
	checkSummaryResults(false, format, reportSummary, chartUri, t)

}

func checkSummaryAll(format SummaryFormat, reportSummary string, chartUri string, t *testing.T) {
	checkSummaryAnnotations(true, format, reportSummary, chartUri, t)
	checkSummaryDigests(true, format, reportSummary, chartUri, t)
	checkSummaryMetadata(true, format, reportSummary, chartUri, t)
	checkSummaryResults(true, format, reportSummary, chartUri, t)

}

func checkSummaryResults(fullReport bool, format SummaryFormat, reportSummary string, chartUri string, t *testing.T) {
	if format == JsonReport {
		if !fullReport {
			require.NotContains(t, reportSummary, "\"annotations\":[{")
			require.NotContains(t, reportSummary, "\"digests\":{")
			require.NotContains(t, reportSummary, "\"metadata\":{")
			require.NotContains(t, reportSummary, "\"name\":\"charts.openshift.io/digest\"")
			require.NotContains(t, reportSummary, "\"chart\":\"sha256")
			require.NotContains(t, reportSummary, "\"chart-uri\":\""+chartUri+"\"")
		}
		require.Contains(t, reportSummary, "\"results\":{")
		require.Contains(t, reportSummary, "\"passed\":")
		require.Contains(t, reportSummary, "\"failed\":")
	} else {
		if !fullReport {
			require.NotContains(t, reportSummary, "digests:")
			require.NotContains(t, reportSummary, "metadata:")
			require.NotContains(t, reportSummary, "- name: charts.openshift.io/digest")
			require.NotContains(t, reportSummary, "chart: sha256:")
			require.NotContains(t, reportSummary, "chart-uri: "+chartUri)
		}
		require.Contains(t, reportSummary, "results:")
		require.Contains(t, reportSummary, "passed:")
		require.Contains(t, reportSummary, "failed:")

	}
}

func checkSummaryAnnotations(fullReport bool, format SummaryFormat, reportSummary string, chartUri string, t *testing.T) {
	if format == JsonReport {
		if !fullReport {
			require.NotContains(t, reportSummary, "\"digests\":{")
			require.NotContains(t, reportSummary, "\"metadata\":{")
			require.NotContains(t, reportSummary, "\"results\":{")
			require.NotContains(t, reportSummary, "\"chart\":\"sha256")
			require.NotContains(t, reportSummary, "\"chart-uri\":\""+chartUri+"\"")
			require.NotContains(t, reportSummary, "\"passed\":")
			require.NotContains(t, reportSummary, "\"failed\":")
		}
		require.Contains(t, reportSummary, "\"annotations\":[{")
		require.Contains(t, reportSummary, "\"name\":\"charts.openshift.io/digest\"")
	} else {
		if !fullReport {
			require.NotContains(t, reportSummary, "digests:")
			require.NotContains(t, reportSummary, "metadata:")
			require.NotContains(t, reportSummary, "results:")
			require.NotContains(t, reportSummary, "chart: sha256:")
			require.NotContains(t, reportSummary, "chart-uri: "+chartUri)
			require.NotContains(t, reportSummary, "passed:")
			require.NotContains(t, reportSummary, "failed:")
		}
		require.Contains(t, reportSummary, "annotations:")
		require.Contains(t, reportSummary, "- name: charts.openshift.io/digest")
	}
}

func checkSummaryMetadata(fullReport bool, format SummaryFormat, reportSummary string, chartUri string, t *testing.T) {
	if format == JsonReport {
		if !fullReport {
			require.NotContains(t, reportSummary, "\"annotations\":[{")
			require.NotContains(t, reportSummary, "\"digests\":{")
			require.NotContains(t, reportSummary, "\"results\":{")
			require.NotContains(t, reportSummary, "\"name\":\"charts.openshift.io/digest\"")
			require.NotContains(t, reportSummary, "\"chart\":\"sha256")
			require.NotContains(t, reportSummary, "\"passed\":")
			require.NotContains(t, reportSummary, "\"failed\":")
		}
		require.Contains(t, reportSummary, "\"metadata\":{")
		require.Contains(t, reportSummary, "\"chart-uri\":\""+chartUri+"\"")
	} else {
		if !fullReport {
			require.NotContains(t, reportSummary, "digests:")
			require.NotContains(t, reportSummary, "results:")
			require.NotContains(t, reportSummary, "- name: charts.openshift.io/digest")
			require.NotContains(t, reportSummary, "chart: sha256:")
			require.NotContains(t, reportSummary, "passed:")
			require.NotContains(t, reportSummary, "failed:")
		}
		require.Contains(t, reportSummary, "metadata:")
		require.Contains(t, reportSummary, "chart-uri: "+chartUri)
	}
}

func checkSummaryDigests(fullReport bool, format SummaryFormat, reportSummary string, chartUri string, t *testing.T) {
	if format == JsonReport {
		if !fullReport {
			require.NotContains(t, reportSummary, "\"annotations\":[{")
			require.NotContains(t, reportSummary, "\"metadata\":{")
			require.NotContains(t, reportSummary, "\"results\":{")
			require.NotContains(t, reportSummary, "\"name\":\"charts.openshift.io/digest\"")
			require.NotContains(t, reportSummary, "\"chart-uri\":\""+chartUri+"\"")
			require.NotContains(t, reportSummary, "\"passed\":")
			require.NotContains(t, reportSummary, "\"failed\":")
		}
		require.Contains(t, reportSummary, "\"digests\":{")
		require.Contains(t, reportSummary, "\"chart\":\"sha256")
	} else {
		if !fullReport {
			require.NotContains(t, reportSummary, "annotations:")
			require.NotContains(t, reportSummary, "metadata:")
			require.NotContains(t, reportSummary, "results:")
			require.NotContains(t, reportSummary, "- name: charts.openshift.io/digest")
			require.NotContains(t, reportSummary, "chart-uri: "+chartUri)
			require.NotContains(t, reportSummary, "passed:")
			require.NotContains(t, reportSummary, "failed:")
		}
		require.Contains(t, reportSummary, "digests:")
		require.Contains(t, reportSummary, "chart: sha256:")
	}
}

func loadChartFromAbsPath(path string) ([]byte, error) {
	// Open the yaml file which defines the tests to run
	reportFile, openErr := os.Open(path)
	if openErr != nil {
		return nil, errors.New(fmt.Sprintf("report path %s: error opening file  %v", path, openErr))
	}

	reportBytes, readErr := ioutil.ReadAll(reportFile)
	if readErr != nil {
		return nil, errors.New(fmt.Sprintf("report path %s: error reading file  %v", path, readErr))
	}

	return reportBytes, nil
}
