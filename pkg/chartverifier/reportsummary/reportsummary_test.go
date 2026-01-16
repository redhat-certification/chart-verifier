package reportsummary

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	apireport "github.com/redhat-certification/chart-verifier/pkg/chartverifier/report"
)

func TestAddStringReport(t *testing.T) {
	chartURI := "https://github.com/redhat-certification/chart-verifier/blob/main/tests/charts/psql-service/0.1.9/psql-service-0.1.9.tgz?raw=true"
	yamlFileReport := "test-reports/report.yaml"
	jsonFileReport := "test-reports/report.json"

	yamlFileReportBytes, yamlReportErr := loadChartFromAbsPath(yamlFileReport)
	require.NoError(t, yamlReportErr)
	jsonFileReportBytes, jsonReportErr := loadChartFromAbsPath(jsonFileReport)
	require.NoError(t, jsonReportErr)

	report, loadErr := apireport.NewReport().SetContent(string(yamlFileReportBytes)).Load()
	require.NoError(t, loadErr)
	reportSummary := NewReportSummary().SetReport(report)
	checkReportSummaries(reportSummary, chartURI, t)

	report, loadErr = apireport.NewReport().SetContent(string(jsonFileReportBytes)).Load()
	require.NoError(t, loadErr)
	reportSummary = NewReportSummary().SetReport(report)
	checkReportSummaries(reportSummary, chartURI, t)
}

func TestBadDigestReport(t *testing.T) {
	yamlFileReport := "test-reports/baddigest/report.yaml"
	jsonFileReport := "test-reports/baddigest/report.json"

	yamlFileReportBytes, yamlReportErr := loadChartFromAbsPath(yamlFileReport)
	require.NoError(t, yamlReportErr)
	jsonFileReportBytes, jsonReportErr := loadChartFromAbsPath(jsonFileReport)
	require.NoError(t, jsonReportErr)

	report, loadErr := apireport.NewReport().SetContent(string(yamlFileReportBytes)).Load()
	require.NoError(t, loadErr)
	_, summaryErr := NewReportSummary().SetReport(report).GetContent(AllSummary, YAMLReport)
	require.Error(t, summaryErr)
	require.Contains(t, fmt.Sprintf("%v", summaryErr), "digest in report did not match report content")

	report, loadErr = apireport.NewReport().SetContent(string(jsonFileReportBytes)).Load()
	require.NoError(t, loadErr)
	_, summaryErr = NewReportSummary().SetReport(report).GetContent(AllSummary, JSONReport)
	require.Error(t, summaryErr)
	require.Contains(t, fmt.Sprintf("%v", summaryErr), "digest in report did not match report content")
}

func TestSkipBadDigestReport(t *testing.T) {
	chartURI := "https://github.com/redhat-certification/chart-verifier/blob/main/tests/charts/psql-service/0.1.9/psql-service-0.1.9.tgz?raw=true"
	yamlFileReport := "test-reports/baddigest/report.yaml"
	jsonFileReport := "test-reports/baddigest/report.json"

	yamlFileReportBytes, yamlReportErr := loadChartFromAbsPath(yamlFileReport)
	require.NoError(t, yamlReportErr)
	jsonFileReportBytes, jsonReportErr := loadChartFromAbsPath(jsonFileReport)
	require.NoError(t, jsonReportErr)

	report, loadErr := apireport.NewReport().SetContent(string(yamlFileReportBytes)).Load()
	require.NoError(t, loadErr)
	reportSummary := NewReportSummary().SetReport(report).SetBoolean(SkipDigestCheck, true)
	checkReportSummaries(reportSummary, chartURI, t)

	report, loadErr = apireport.NewReport().SetContent(string(jsonFileReportBytes)).Load()
	require.NoError(t, loadErr)
	reportSummary = NewReportSummary().SetReport(report).SetBoolean(SkipDigestCheck, true)
	checkReportSummaries(reportSummary, chartURI, t)
}

func TestMissingDigestReport(t *testing.T) {
	yamlFileReport := "test-reports/missingdigest/report.yaml"
	jsonFileReport := "test-reports/missingdigest/report.json"

	yamlFileReportBytes, yamlReportErr := loadChartFromAbsPath(yamlFileReport)
	require.NoError(t, yamlReportErr)
	jsonFileReportBytes, jsonReportErr := loadChartFromAbsPath(jsonFileReport)
	require.NoError(t, jsonReportErr)

	report, loadErr := apireport.NewReport().SetContent(string(yamlFileReportBytes)).Load()
	require.NoError(t, loadErr)
	_, summaryErr := NewReportSummary().SetReport(report).GetContent(AllSummary, YAMLReport)
	require.Error(t, summaryErr)
	require.Contains(t, fmt.Sprintf("%v", summaryErr), "report does not contain expected report digest")

	report, loadErr = apireport.NewReport().SetContent(string(jsonFileReportBytes)).Load()
	require.NoError(t, loadErr)
	_, summaryErr = NewReportSummary().SetReport(report).GetContent(AllSummary, JSONReport)
	require.Error(t, summaryErr)
	require.Contains(t, fmt.Sprintf("%v", summaryErr), "report does not contain expected report digest")
}

func TestPreDigestReport(t *testing.T) {
	chartURI := "https://github.com/redhat-certification/chart-verifier/blob/main/tests/charts/psql-service/0.1.9/psql-service-0.1.9.tgz?raw=true"
	yamlFileReport := "test-reports/predigest/report.yaml"
	jsonFileReport := "test-reports/predigest/report.json"

	yamlFileReportBytes, yamlReportErr := loadChartFromAbsPath(yamlFileReport)
	require.NoError(t, yamlReportErr)
	jsonFileReportBytes, jsonReportErr := loadChartFromAbsPath(jsonFileReport)
	require.NoError(t, jsonReportErr)

	report, loadErr := apireport.NewReport().SetContent(string(yamlFileReportBytes)).Load()
	require.NoError(t, loadErr)
	reportSummary := NewReportSummary().SetReport(report)
	checkReportSummaries(reportSummary, chartURI, t)

	report, loadErr = apireport.NewReport().SetContent(string(jsonFileReportBytes)).Load()
	require.NoError(t, loadErr)
	reportSummary = NewReportSummary().SetReport(report)
	checkReportSummaries(reportSummary, chartURI, t)
}

func TestProviderDeliveryReport(t *testing.T) {
	chartURI := "N/A"
	yamlFileReport := "test-reports/providerdelivery/report.yaml"
	jsonFileReport := "test-reports/providerdelivery/report.json"

	yamlFileReportBytes, yamlReportErr := loadChartFromAbsPath(yamlFileReport)
	require.NoError(t, yamlReportErr)
	jsonFileReportBytes, jsonReportErr := loadChartFromAbsPath(jsonFileReport)
	require.NoError(t, jsonReportErr)

	report, loadErr := apireport.NewReport().SetContent(string(yamlFileReportBytes)).Load()
	require.NoError(t, loadErr)
	reportSummary := NewReportSummary().SetReport(report)
	checkReportSummaries(reportSummary, chartURI, t)

	report, loadErr = apireport.NewReport().SetContent(string(jsonFileReportBytes)).Load()
	require.NoError(t, loadErr)
	reportSummary = NewReportSummary().SetReport(report)
	checkReportSummaries(reportSummary, chartURI, t)
}

func TestAddUrlReport(t *testing.T) {
	yamlURLReport := "https://github.com/redhat-certification/chart-verifier/blob/main/cmd/test/report.yaml?raw=true"
	url, loadURLErr := url.Parse(yamlURLReport)
	require.NoError(t, loadURLErr)

	chartURI := "internal/chartverifier/checks/chart-0.1.0-v3.valid.tgz"
	report, loadErr := apireport.NewReport().SetURL(url).Load()
	require.NoError(t, loadErr)
	reportSummary := NewReportSummary().SetReport(report)

	checkReportSummaries(reportSummary, chartURI, t)
}

func checkReportSummaries(summary APIReportSummary, chartURI string, t *testing.T) {
	checkReportSummariesFormat(YAMLReport, summary, chartURI, t)
	checkReportSummariesFormat(JSONReport, summary, chartURI, t)
}

func checkReportSummariesFormat(format SummaryFormat, summary APIReportSummary, chartURI string, t *testing.T) {
	reportSummary, reportSummaryErr := summary.GetContent(AllSummary, format)
	require.NoError(t, reportSummaryErr)
	checkSummaryAll(format, reportSummary, chartURI, t)

	reportSummary, reportSummaryErr = summary.GetContent(DigestsSummary, format)
	require.NoError(t, reportSummaryErr)
	checkSummaryDigests(false, format, reportSummary, chartURI, t)

	reportSummary, reportSummaryErr = summary.GetContent(MetadataSummary, format)
	require.NoError(t, reportSummaryErr)
	checkSummaryMetadata(false, format, reportSummary, chartURI, t)

	reportSummary, reportSummaryErr = summary.GetContent(AnnotationsSummary, format)
	require.NoError(t, reportSummaryErr)
	checkSummaryAnnotations(false, format, reportSummary, chartURI, t)

	reportSummary, reportSummaryErr = summary.GetContent(ResultsSummary, format)
	require.NoError(t, reportSummaryErr)
	checkSummaryResults(false, format, reportSummary, chartURI, t)
}

func checkSummaryAll(format SummaryFormat, reportSummary string, chartURI string, t *testing.T) {
	checkSummaryAnnotations(true, format, reportSummary, chartURI, t)
	checkSummaryDigests(true, format, reportSummary, chartURI, t)
	checkSummaryMetadata(true, format, reportSummary, chartURI, t)
	checkSummaryResults(true, format, reportSummary, chartURI, t)
}

func checkSummaryResults(fullReport bool, format SummaryFormat, reportSummary string, chartURI string, t *testing.T) {
	if format == JSONReport {
		if !fullReport {
			require.NotContains(t, reportSummary, "\"annotations\":[{")
			require.NotContains(t, reportSummary, "\"digests\":{")
			require.NotContains(t, reportSummary, "\"metadata\":{")
			require.NotContains(t, reportSummary, "\"name\":\"charts.openshift.io/digest\"")
			require.NotContains(t, reportSummary, "\"chart\":\"sha256")
			require.NotContains(t, reportSummary, "\"chart-uri\":\""+chartURI+"\"")
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
			require.NotContains(t, reportSummary, "chart-uri: "+chartURI)
		}
		require.Contains(t, reportSummary, "results:")
		require.Contains(t, reportSummary, "passed:")
		require.Contains(t, reportSummary, "failed:")
	}
}

func checkSummaryAnnotations(fullReport bool, format SummaryFormat, reportSummary string, chartURI string, t *testing.T) {
	if format == JSONReport {
		if !fullReport {
			require.NotContains(t, reportSummary, "\"digests\":{")
			require.NotContains(t, reportSummary, "\"metadata\":{")
			require.NotContains(t, reportSummary, "\"results\":{")
			require.NotContains(t, reportSummary, "\"chart\":\"sha256")
			require.NotContains(t, reportSummary, "\"chart-uri\":\""+chartURI+"\"")
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
			require.NotContains(t, reportSummary, "chart-uri: "+chartURI)
			require.NotContains(t, reportSummary, "passed:")
			require.NotContains(t, reportSummary, "failed:")
		}
		require.Contains(t, reportSummary, "annotations:")
		require.Contains(t, reportSummary, "- name: charts.openshift.io/digest")
	}
}

func checkSummaryMetadata(fullReport bool, format SummaryFormat, reportSummary string, chartURI string, t *testing.T) {
	if format == JSONReport {
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
		require.Contains(t, reportSummary, "\"chart-uri\":\""+chartURI+"\"")
		require.Contains(t, reportSummary, "\"webCatalogOnly\":")
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
		require.Contains(t, reportSummary, "chart-uri: "+chartURI)
	}
}

func checkSummaryDigests(fullReport bool, format SummaryFormat, reportSummary string, chartURI string, t *testing.T) {
	if format == JSONReport {
		if !fullReport {
			require.NotContains(t, reportSummary, "\"annotations\":[{")
			require.NotContains(t, reportSummary, "\"metadata\":{")
			require.NotContains(t, reportSummary, "\"results\":{")
			require.NotContains(t, reportSummary, "\"name\":\"charts.openshift.io/digest\"")
			require.NotContains(t, reportSummary, "\"chart-uri\":\""+chartURI+"\"")
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
			require.NotContains(t, reportSummary, "chart-uri: "+chartURI)
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
		return nil, fmt.Errorf("report path %s: error opening file  %v", path, openErr)
	}
	defer reportFile.Close()

	reportBytes, readErr := io.ReadAll(reportFile)
	if readErr != nil {
		return nil, fmt.Errorf("report path %s: error reading file  %v", path, readErr)
	}

	return reportBytes, nil
}
