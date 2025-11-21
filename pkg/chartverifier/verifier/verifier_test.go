package verifier

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/redhat-certification/chart-verifier/internal/tool"
	apichecks "github.com/redhat-certification/chart-verifier/pkg/chartverifier/checks"
	apireport "github.com/redhat-certification/chart-verifier/pkg/chartverifier/report"
	apireportsummary "github.com/redhat-certification/chart-verifier/pkg/chartverifier/reportsummary"
	apiversion "github.com/redhat-certification/chart-verifier/pkg/chartverifier/version"
)

func TestVerifyApi(t *testing.T) {
	chartURI := "../../../internal/chartverifier/checks/chart-0.1.0-v3.valid.tgz"
	verifier, reportErr := NewVerifier().
		UnEnableChecks([]apichecks.CheckName{apichecks.ChartTesting, apichecks.ImagesAreCertified, apichecks.ClusterIsNotEOL}).
		Run(chartURI)

	require.NoError(t, reportErr)

	reportSummary := apireportsummary.NewReportSummary().SetReport(verifier.GetReport())
	checkReportSummaries(reportSummary, chartURI, t)
	require.True(t, verifier.GetReport().Metadata.ToolMetadata.Version == apiversion.GetVersion())
}

func TestProfiles(t *testing.T) {
	commandSet := make(map[string]interface{})
	commandSet["profile.vendortype"] = "redhat"

	verifier, RunErr := NewVerifier().
		SetValues(CommandSet, commandSet).
		UnEnableChecks([]apichecks.CheckName{apichecks.ChartTesting, apichecks.ImagesAreCertified, apichecks.ClusterIsNotEOL}).
		Run("../../../internal/chartverifier/checks/chart-0.1.0-v3.valid.tgz")
	require.NoError(t, RunErr)

	report, reportErr := verifier.GetReport().GetContent(apireport.YamlReport)
	require.NoError(t, reportErr)
	require.Contains(t, report, "VendorType: redhat")
}

func TestProfilesDeveloperConsole(t *testing.T) {
	commandSet := make(map[string]interface{})
	commandSet["profile.vendortype"] = "developer-console"

	verifier, RunErr := NewVerifier().
		SetValues(CommandSet, commandSet).
		Run("../../../internal/chartverifier/checks/chart-0.1.0-v3.valid.tgz")
	require.NoError(t, RunErr)

	report, reportErr := verifier.GetReport().GetContent(apireport.YamlReport)
	require.NoError(t, reportErr)
	require.Contains(t, report, "VendorType: developer-console")
	require.NotContains(t, report, "check: v1.0/chart-testing")
	require.NotContains(t, report, "check: v1.0/images-are-certified")
}

func TestWebCatalogOnly(t *testing.T) {
	commandSet := make(map[string]interface{})
	commandSet["profile.vendortype"] = "redhat"

	verifier, RunErr := NewVerifier().
		SetBoolean(WebCatalogOnly, true).
		UnEnableChecks([]apichecks.CheckName{apichecks.ChartTesting, apichecks.ImagesAreCertified, apichecks.ClusterIsNotEOL}).
		Run("../../../internal/chartverifier/checks/chart-0.1.0-v3.valid.tgz")
	require.NoError(t, RunErr)

	reportContent, reportErr := verifier.GetReport().GetContent(apireport.YamlReport)
	require.NoError(t, reportErr)
	require.Contains(t, reportContent, "webCatalogOnly: true")
	require.Contains(t, reportContent, "chart-uri:")
	require.NotContains(t, reportContent, "chart-0.1.0-v3.valid.tgz")

	report := verifier.GetReport()
	require.True(t, report.Metadata.ToolMetadata.WebCatalogOnly)
	require.NotContains(t, report.Metadata.ToolMetadata.ChartUri, "chart-0.1.0-v3.valid.tgz")

	verifier, RunErr = NewVerifier().
		SetBoolean(WebCatalogOnly, false).
		UnEnableChecks([]apichecks.CheckName{apichecks.ChartTesting, apichecks.ImagesAreCertified, apichecks.ClusterIsNotEOL}).
		Run("../../../internal/chartverifier/checks/chart-0.1.0-v3.valid.tgz")
	require.NoError(t, RunErr)

	reportContent, reportErr = verifier.GetReport().GetContent(apireport.YamlReport)
	require.NoError(t, reportErr)
	require.Contains(t, reportContent, "webCatalogOnly: false")
	require.Contains(t, reportContent, "chart-uri:")
	require.Contains(t, reportContent, "chart-0.1.0-v3.valid.tgz")

	report = verifier.GetReport()
	require.False(t, report.Metadata.ToolMetadata.WebCatalogOnly)
	require.Contains(t, report.Metadata.ToolMetadata.ChartUri, "chart-0.1.0-v3.valid.tgz")
}

func TestBadFlags(t *testing.T) {
	_, runErr := NewVerifier().
		SetString(StringKey("badStringKey"), []string{"Bad key value"}).
		Run("../../../internal/chartverifier/checks/chart-0.1.0-v3.valid.tgz")
	require.Error(t, runErr)
	require.Contains(t, fmt.Sprint(runErr), "invalid string key name: badStringKey")

	_, runErr = NewVerifier().
		UnEnableChecks([]apichecks.CheckName{apichecks.CheckName("Bad-Check")}).
		Run("../../../internal/chartverifier/checks/chart-0.1.0-v3.valid.tgz")
	require.Error(t, runErr)
	require.Contains(t, fmt.Sprint(runErr), "invalid check name : Bad-Check")

	badValueSet := make(map[string]interface{})
	badValueSet["key"] = "value"
	_, runErr = NewVerifier().
		SetValues(ValuesKey("BadValueKey"), badValueSet).
		Run("../../../internal/chartverifier/checks/chart-0.1.0-v3.valid.tgz")
	require.Error(t, runErr)
	require.Contains(t, fmt.Sprint(runErr), "invalid values key name: BadValueKey")

	_, runErr = NewVerifier().
		SetBoolean(BooleanKey("BadBooleanKey"), false).
		Run("../../../internal/chartverifier/checks/chart-0.1.0-v3.valid.tgz")
	require.Error(t, runErr)
	require.Contains(t, fmt.Sprint(runErr), "invalid boolean key name: BadBooleanKey")

	_, runErr = NewVerifier().
		SetDuration(DurationKey("BadDurationKey"), 3000000).
		Run("../../../internal/chartverifier/checks/chart-0.1.0-v3.valid.tgz")
	require.Error(t, runErr)
	require.Contains(t, fmt.Sprint(runErr), "invalid duration key name: BadDurationKey")

	_, runErr = NewVerifier().
		Run("")
	require.Error(t, runErr)
	require.Contains(t, fmt.Sprint(runErr), "run error: chart_uri is required")
}

func TestSignedChart(t *testing.T) {
	chartURI := "../../../tests/charts/psql-service/0.1.11/psql-service-0.1.11.tgz"
	key1 := "../../../tests/charts/psql-service/0.1.11/psql-service-0.1.11.tgz.key"
	encodedKey1, err1 := tool.GetEncodedKey(key1)
	require.NoError(t, err1)

	key2 := "../../../tests/charts/psql-service/0.1.11/psql-service-0.1.11.tgz.badkey"
	encodedKey2, err2 := tool.GetEncodedKey(key2)
	require.NoError(t, err2)

	// good key first
	verifier, reportErr := NewVerifier().
		UnEnableChecks([]apichecks.CheckName{apichecks.ChartTesting, apichecks.ImagesAreCertified, apichecks.ClusterIsNotEOL}).
		SetString(PGPPublicKey, []string{encodedKey1, encodedKey2}).
		Run(chartURI)

	require.NoError(t, reportErr)
	report := verifier.GetReport()

	sigResultFound := false
	for _, checkResult := range report.Results {
		if strings.Contains((string)(checkResult.Check), (string)(apichecks.SignatureIsValid)) {
			require.Equal(t, checkResult.Outcome, apireport.PassOutcomeType)
			require.Contains(t, checkResult.Reason, "Chart is signed")
			require.Contains(t, checkResult.Reason, "Signature verification passed")
			sigResultFound = true
		}
	}
	require.True(t, sigResultFound)

	// bad key first
	verifier, reportErr = NewVerifier().
		UnEnableChecks([]apichecks.CheckName{apichecks.ChartTesting, apichecks.ImagesAreCertified, apichecks.ClusterIsNotEOL}).
		SetString(PGPPublicKey, []string{encodedKey2, encodedKey1}).
		Run(chartURI)

	require.NoError(t, reportErr)
	report = verifier.GetReport()

	sigResultFound = false
	for _, checkResult := range report.Results {
		if strings.Contains((string)(checkResult.Check), (string)(apichecks.SignatureIsValid)) {
			require.Equal(t, checkResult.Outcome, apireport.PassOutcomeType)
			require.Contains(t, checkResult.Reason, "Chart is signed")
			require.Contains(t, checkResult.Reason, "Signature verification passed")
			sigResultFound = true
		}
	}
	require.True(t, sigResultFound)
}

func checkReportSummaries(summary apireportsummary.APIReportSummary, chartURI string, t *testing.T) {
	checkReportSummariesFormat(apireportsummary.YAMLReport, summary, chartURI, t)
	checkReportSummariesFormat(apireportsummary.JSONReport, summary, chartURI, t)
}

func checkReportSummariesFormat(format apireportsummary.SummaryFormat, summary apireportsummary.APIReportSummary, chartURI string, t *testing.T) {
	reportSummary, reportSummaryErr := summary.GetContent(apireportsummary.AllSummary, format)
	require.NoError(t, reportSummaryErr)
	checkSummaryAll(format, reportSummary, chartURI, t)

	reportSummary, reportSummaryErr = summary.GetContent(apireportsummary.DigestsSummary, format)
	require.NoError(t, reportSummaryErr)
	checkSummaryDigests(false, format, reportSummary, chartURI, t)

	reportSummary, reportSummaryErr = summary.GetContent(apireportsummary.MetadataSummary, format)
	require.NoError(t, reportSummaryErr)
	checkSummaryMetadata(false, format, reportSummary, chartURI, t)

	reportSummary, reportSummaryErr = summary.GetContent(apireportsummary.AnnotationsSummary, format)
	require.NoError(t, reportSummaryErr)
	checkSummaryAnnotations(false, format, reportSummary, chartURI, t)

	reportSummary, reportSummaryErr = summary.GetContent(apireportsummary.ResultsSummary, format)
	require.NoError(t, reportSummaryErr)
	checkSummaryResults(false, format, reportSummary, chartURI, t)
}

func checkSummaryAll(format apireportsummary.SummaryFormat, reportSummary string, chartURI string, t *testing.T) {
	checkSummaryAnnotations(true, format, reportSummary, chartURI, t)
	checkSummaryDigests(true, format, reportSummary, chartURI, t)
	checkSummaryMetadata(true, format, reportSummary, chartURI, t)
	checkSummaryResults(true, format, reportSummary, chartURI, t)
}

func checkSummaryResults(fullReport bool, format apireportsummary.SummaryFormat, reportSummary string, chartURI string, t *testing.T) {
	if format == apireportsummary.JSONReport {
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

func checkSummaryAnnotations(fullReport bool, format apireportsummary.SummaryFormat, reportSummary string, chartURI string, t *testing.T) {
	if format == apireportsummary.JSONReport {
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

func checkSummaryMetadata(fullReport bool, format apireportsummary.SummaryFormat, reportSummary string, chartURI string, t *testing.T) {
	if format == apireportsummary.JSONReport {
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

func checkSummaryDigests(fullReport bool, format apireportsummary.SummaryFormat, reportSummary string, chartURI string, t *testing.T) {
	if format == apireportsummary.JSONReport {
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
