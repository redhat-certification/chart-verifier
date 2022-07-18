package verifier

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	apichecks "github.com/redhat-certification/chart-verifier/pkg/chartverifier/checks"
	apireport "github.com/redhat-certification/chart-verifier/pkg/chartverifier/report"
	apireportsummary "github.com/redhat-certification/chart-verifier/pkg/chartverifier/reportsummary"
)

func TestVerifyApi(t *testing.T) {

	chartUri := "../../../internal/chartverifier/checks/chart-0.1.0-v3.valid.tgz"
	verifier, reportErr := NewVerifier().
		UnEnableChecks([]apichecks.CheckName{apichecks.ChartTesting, apichecks.ImagesAreCertified}).
		Run(chartUri)

	require.NoError(t, reportErr)

	reportSummary := apireportsummary.NewReportSummary().SetReport(verifier.GetReport())
	checkReportSummaries(reportSummary, chartUri, t)

}

func TestProfiles(t *testing.T) {

	commandSet := make(map[string]interface{})
	commandSet["profile.vendortype"] = "redhat"

	verifier, RunErr := NewVerifier().
		SetValues(CommandSet, commandSet).
		UnEnableChecks([]apichecks.CheckName{apichecks.ChartTesting, apichecks.ImagesAreCertified}).
		Run("../../../internal/chartverifier/checks/chart-0.1.0-v3.valid.tgz")
	require.NoError(t, RunErr)

	report, reportErr := verifier.GetReport().GetContent(apireport.YamlReport)
	require.NoError(t, reportErr)
	require.Contains(t, report, "VendorType: redhat")

}

func TestProviderDelivery(t *testing.T) {

	commandSet := make(map[string]interface{})
	commandSet["profile.vendortype"] = "redhat"

	verifier, RunErr := NewVerifier().
		SetBoolean(ProviderDelivery, true).
		UnEnableChecks([]apichecks.CheckName{apichecks.ChartTesting, apichecks.ImagesAreCertified}).
		Run("../../../internal/chartverifier/checks/chart-0.1.0-v3.valid.tgz")
	require.NoError(t, RunErr)

	reportContent, reportErr := verifier.GetReport().GetContent(apireport.YamlReport)
	require.NoError(t, reportErr)
	require.Contains(t, reportContent, "providerControlledDelivery: true")
	require.Contains(t, reportContent, "chart-uri:")
	require.NotContains(t, reportContent, "chart-0.1.0-v3.valid.tgz")

	report := verifier.GetReport()
	require.True(t, report.Metadata.ToolMetadata.ProviderDelivery)
	require.NotContains(t, report.Metadata.ToolMetadata.ChartUri, "chart-0.1.0-v3.valid.tgz")

	verifier, RunErr = NewVerifier().
		SetBoolean(ProviderDelivery, false).
		UnEnableChecks([]apichecks.CheckName{apichecks.ChartTesting, apichecks.ImagesAreCertified}).
		Run("../../../internal/chartverifier/checks/chart-0.1.0-v3.valid.tgz")
	require.NoError(t, RunErr)

	reportContent, reportErr = verifier.GetReport().GetContent(apireport.YamlReport)
	require.NoError(t, reportErr)
	require.Contains(t, reportContent, "providerControlledDelivery: false")
	require.Contains(t, reportContent, "chart-uri:")
	require.Contains(t, reportContent, "chart-0.1.0-v3.valid.tgz")

	report = verifier.GetReport()
	require.False(t, report.Metadata.ToolMetadata.ProviderDelivery)
	require.Contains(t, report.Metadata.ToolMetadata.ChartUri, "chart-0.1.0-v3.valid.tgz")

}

func TestBadFlags(t *testing.T) {

	_, runErr := NewVerifier().
		SetString(StringKey("badStringKey"), []string{"Bad key value"}).
		Run("../../../internal/chartverifier/checks/chart-0.1.0-v3.valid.tgz")
	require.Error(t, runErr)
	require.Contains(t, fmt.Sprint(runErr), "Invalid string key name: badStringKey")

	_, runErr = NewVerifier().
		UnEnableChecks([]apichecks.CheckName{apichecks.CheckName("Bad-Check")}).
		Run("../../../internal/chartverifier/checks/chart-0.1.0-v3.valid.tgz")
	require.Error(t, runErr)
	require.Contains(t, fmt.Sprint(runErr), "Invalid check name : Bad-Check")

	badValueSet := make(map[string]interface{})
	badValueSet["key"] = "value"
	_, runErr = NewVerifier().
		SetValues(ValuesKey("BadValueKey"), badValueSet).
		Run("../../../internal/chartverifier/checks/chart-0.1.0-v3.valid.tgz")
	require.Error(t, runErr)
	require.Contains(t, fmt.Sprint(runErr), "Invalid values key name: BadValueKey")

	_, runErr = NewVerifier().
		SetBoolean(BooleanKey("BadBooleanKey"), false).
		Run("../../../internal/chartverifier/checks/chart-0.1.0-v3.valid.tgz")
	require.Error(t, runErr)
	require.Contains(t, fmt.Sprint(runErr), "Invalid boolean key name: BadBooleanKey")

	_, runErr = NewVerifier().
		SetDuration(DurationKey("BadDurationKey"), 3000000).
		Run("../../../internal/chartverifier/checks/chart-0.1.0-v3.valid.tgz")
	require.Error(t, runErr)
	require.Contains(t, fmt.Sprint(runErr), "Invalid duration key name: BadDurationKey")

	_, runErr = NewVerifier().
		Run("")
	require.Error(t, runErr)
	require.Contains(t, fmt.Sprint(runErr), "run error: chart_uri is required")

}

func checkReportSummaries(summary apireportsummary.APIReportSummary, chartUri string, t *testing.T) {
	checkReportSummariesFormat(apireportsummary.YamlReport, summary, chartUri, t)
	checkReportSummariesFormat(apireportsummary.JsonReport, summary, chartUri, t)
}

func checkReportSummariesFormat(format apireportsummary.SummaryFormat, summary apireportsummary.APIReportSummary, chartUri string, t *testing.T) {

	reportSummary, reportSummaryErr := summary.GetContent(apireportsummary.AllSummary, format)
	require.NoError(t, reportSummaryErr)
	checkSummaryAll(format, reportSummary, chartUri, t)

	reportSummary, reportSummaryErr = summary.GetContent(apireportsummary.DigestsSummary, format)
	require.NoError(t, reportSummaryErr)
	checkSummaryDigests(false, format, reportSummary, chartUri, t)

	reportSummary, reportSummaryErr = summary.GetContent(apireportsummary.MetadataSummary, format)
	require.NoError(t, reportSummaryErr)
	checkSummaryMetadata(false, format, reportSummary, chartUri, t)

	reportSummary, reportSummaryErr = summary.GetContent(apireportsummary.AnnotationsSummary, format)
	require.NoError(t, reportSummaryErr)
	checkSummaryAnnotations(false, format, reportSummary, chartUri, t)

	reportSummary, reportSummaryErr = summary.GetContent(apireportsummary.ResultsSummary, format)
	require.NoError(t, reportSummaryErr)
	checkSummaryResults(false, format, reportSummary, chartUri, t)

}

func checkSummaryAll(format apireportsummary.SummaryFormat, reportSummary string, chartUri string, t *testing.T) {
	checkSummaryAnnotations(true, format, reportSummary, chartUri, t)
	checkSummaryDigests(true, format, reportSummary, chartUri, t)
	checkSummaryMetadata(true, format, reportSummary, chartUri, t)
	checkSummaryResults(true, format, reportSummary, chartUri, t)

}

func checkSummaryResults(fullReport bool, format apireportsummary.SummaryFormat, reportSummary string, chartUri string, t *testing.T) {
	if format == apireportsummary.JsonReport {
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

func checkSummaryAnnotations(fullReport bool, format apireportsummary.SummaryFormat, reportSummary string, chartUri string, t *testing.T) {
	if format == apireportsummary.JsonReport {
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

func checkSummaryMetadata(fullReport bool, format apireportsummary.SummaryFormat, reportSummary string, chartUri string, t *testing.T) {
	if format == apireportsummary.JsonReport {
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

func checkSummaryDigests(fullReport bool, format apireportsummary.SummaryFormat, reportSummary string, chartUri string, t *testing.T) {
	if format == apireportsummary.JsonReport {
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
