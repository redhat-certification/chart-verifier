package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/redhat-certification/chart-verifier/pkg/chartverifier/profiles"
	"github.com/redhat-certification/chart-verifier/pkg/chartverifier/report"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	helmchart "helm.sh/helm/v3/pkg/chart"
)

func TestReport(t *testing.T) {

	var expectedAnnotations []report.Annotation
	annotation1 := report.Annotation{Name: fmt.Sprintf("%s/%s", report.DefaultAnnotationsPrefix, report.DigestsAnnotationName), Value: "sha256:0c1c44def5c5de45212d90396062e18e0311b07789f477268fbf233c1783dbd0"}
	annotation2 := report.Annotation{Name: fmt.Sprintf("%s/%s", report.DefaultAnnotationsPrefix, report.CertifiedOCPVersionAnnotationName), Value: "4.7.8"}
	annotation3 := report.Annotation{Name: fmt.Sprintf("%s/%s", report.DefaultAnnotationsPrefix, report.LastCertifiedTimestampAnnotationName), Value: "2021-07-06T10:28:01.09604-04:00"}
	expectedAnnotations = append(expectedAnnotations, annotation1, annotation2, annotation3)

	expectedResults := &report.ResultsReport{}
	expectedResults.Passed = "11"
	expectedResults.Failed = "1"

	expectedMetadata := &report.MetadataReport{}
	expectedMetadata.ProfileVersion = "v1.1"
	expectedMetadata.ProfileVendorType = "redhat"
	expectedMetadata.ChartUri = "pkg/chartverifier/checks/chart-0.1.0-v3.valid.tgz"
	expectedMetadata.Chart = &helmchart.Metadata{Name: "chart", Version: "0.1.0-v3.valid"}

	expectedDigests := &report.DigestReport{}
	expectedDigests.PackageDigest = "4f29f2a95bf2b9a1c62fd215b079a01bdc5a38e9b4ff874d0fa21d0afca2e76d"
	expectedDigests.ChartDigest = "sha256:0c1c44def5c5de45212d90396062e18e0311b07789f477268fbf233c1783dbd0"

	t.Run("Should fail when no argument is given", func(t *testing.T) {
		cmd := NewReportCmd(viper.New())
		outBuf := bytes.NewBufferString("")
		cmd.SetOut(outBuf)
		errBuf := bytes.NewBufferString("")
		cmd.SetErr(errBuf)

		require.Error(t, cmd.Execute())
	})

	t.Run("Should fail when one argument is given", func(t *testing.T) {
		cmd := NewReportCmd(viper.New())
		outBuf := bytes.NewBufferString("")
		cmd.SetOut(outBuf)
		errBuf := bytes.NewBufferString("")
		cmd.SetErr(errBuf)

		cmd.SetArgs([]string{
			"test/report.yaml",
		})

		require.Error(t, cmd.Execute())
	})

	t.Run("Should fail when bad subcommand is given", func(t *testing.T) {
		cmd := NewReportCmd(viper.New())
		outBuf := bytes.NewBufferString("")
		cmd.SetOut(outBuf)
		errBuf := bytes.NewBufferString("")
		cmd.SetErr(errBuf)

		cmd.SetArgs([]string{
			"None",
			"test/report.yaml",
		})
		require.Error(t, cmd.Execute())
	})

	t.Run("Should pass for subcommand annotations", func(t *testing.T) {
		cmd := NewReportCmd(viper.New())
		outBuf := bytes.NewBufferString("")
		cmd.SetOut(outBuf)
		errBuf := bytes.NewBufferString("")
		cmd.SetErr(errBuf)

		cmd.SetArgs([]string{
			report.AnnotationsCommandName,
			"test/report.yaml",
		})
		require.NoError(t, cmd.Execute())

		testReport := report.OutputReport{}
		require.NoError(t, json.Unmarshal([]byte(outBuf.String()), &testReport))
		require.True(t, compareAnnotations(expectedAnnotations, testReport.AnnotationsReport))

	})

	t.Run("Should pass for subcommand results", func(t *testing.T) {
		cmd := NewReportCmd(viper.New())
		outBuf := bytes.NewBufferString("")
		cmd.SetOut(outBuf)
		errBuf := bytes.NewBufferString("")
		cmd.SetErr(errBuf)

		cmd.SetArgs([]string{
			report.ResultsCommandName,
			"test/report.yaml",
		})
		require.NoError(t, cmd.Execute())

		testReport := report.OutputReport{}
		require.NoError(t, json.Unmarshal([]byte(outBuf.String()), &testReport))
		require.True(t, compareResults(expectedResults, testReport.ResultsReport))
	})

	t.Run("Should pass for subcommand metadata", func(t *testing.T) {
		cmd := NewReportCmd(viper.New())
		outBuf := bytes.NewBufferString("")
		cmd.SetOut(outBuf)
		errBuf := bytes.NewBufferString("")
		cmd.SetErr(errBuf)

		cmd.SetArgs([]string{
			report.MetadataCommandName,
			"test/report.yaml",
		})
		require.NoError(t, cmd.Execute())

		testReport := report.OutputReport{}
		require.NoError(t, json.Unmarshal([]byte(outBuf.String()), &testReport))
		require.True(t, compareMetadata(expectedMetadata, testReport.MetadataReport))
	})

	t.Run("Should pass for subcommand digests", func(t *testing.T) {
		cmd := NewReportCmd(viper.New())
		outBuf := bytes.NewBufferString("")
		cmd.SetOut(outBuf)
		errBuf := bytes.NewBufferString("")
		cmd.SetErr(errBuf)

		cmd.SetArgs([]string{
			report.DigestsCommandName,
			"test/report.yaml",
		})
		require.NoError(t, cmd.Execute())

		testReport := report.OutputReport{}
		require.NoError(t, json.Unmarshal([]byte(outBuf.String()), &testReport))
		require.True(t, compareDigests(expectedDigests, testReport.DigestsReport))
	})

	t.Run("Should pass for subcommand all", func(t *testing.T) {
		cmd := NewReportCmd(viper.New())
		outBuf := bytes.NewBufferString("")
		cmd.SetOut(outBuf)
		errBuf := bytes.NewBufferString("")
		cmd.SetErr(errBuf)

		cmd.SetArgs([]string{
			report.AllCommandsName,
			"test/report.yaml",
		})
		require.NoError(t, cmd.Execute())

		testReport := report.OutputReport{}
		require.NoError(t, json.Unmarshal([]byte(outBuf.String()), &testReport))
		require.True(t, compareAnnotations(expectedAnnotations, testReport.AnnotationsReport))
		require.True(t, compareDigests(expectedDigests, testReport.DigestsReport))
		require.True(t, compareMetadata(expectedMetadata, testReport.MetadataReport))
		require.True(t, compareResults(expectedResults, testReport.ResultsReport))
	})

	t.Run("Should pass for annotation prefix", func(t *testing.T) {
		cmd := NewReportCmd(viper.New())
		outBuf := bytes.NewBufferString("")
		cmd.SetOut(outBuf)
		errBuf := bytes.NewBufferString("")
		cmd.SetErr(errBuf)

		annotationPrefix := "charts.testing.io"
		cmd.SetArgs([]string{
			"--set", fmt.Sprintf("%s=%s", report.AnnotationsPrefixConfigName, annotationPrefix),
			report.AnnotationsCommandName,
			"test/report.yaml",
		})
		require.NoError(t, cmd.Execute())

		var expectedPrefixAnnotations []report.Annotation
		annotationP1 := report.Annotation{Name: fmt.Sprintf("%s/%s", annotationPrefix, report.DigestsAnnotationName), Value: "sha256:0c1c44def5c5de45212d90396062e18e0311b07789f477268fbf233c1783dbd0"}
		annotationP2 := report.Annotation{Name: fmt.Sprintf("%s/%s", annotationPrefix, report.CertifiedOCPVersionAnnotationName), Value: "4.7.8"}
		annotationP3 := report.Annotation{Name: fmt.Sprintf("%s/%s", annotationPrefix, report.LastCertifiedTimestampAnnotationName), Value: "2021-07-06T10:28:01.09604-04:00"}
		expectedPrefixAnnotations = append(expectedPrefixAnnotations, annotationP1, annotationP2, annotationP3)

		testReport := report.OutputReport{}
		require.NoError(t, json.Unmarshal([]byte(outBuf.String()), &testReport))
		require.True(t, compareAnnotations(expectedPrefixAnnotations, testReport.AnnotationsReport))
	})

	t.Run("Should pass for community profile", func(t *testing.T) {
		cmd := NewReportCmd(viper.New())
		outBuf := bytes.NewBufferString("")
		cmd.SetOut(outBuf)
		errBuf := bytes.NewBufferString("")
		cmd.SetErr(errBuf)

		cmd.SetArgs([]string{
			"--set", fmt.Sprintf("%s=%s", profiles.VendorTypeConfigName, "community"),
			report.AllCommandsName,
			"test/report.yaml",
		})
		require.NoError(t, cmd.Execute())

		expectedCommunityResults := &report.ResultsReport{}
		expectedCommunityResults.Passed = "1"
		expectedCommunityResults.Failed = "0"

		testReport := report.OutputReport{}
		require.NoError(t, json.Unmarshal([]byte(outBuf.String()), &testReport))

		require.True(t, compareMetadata(expectedMetadata, testReport.MetadataReport))
		require.True(t, compareResults(expectedCommunityResults, testReport.ResultsReport))

	})

	t.Run("Should pass for invalid profile version", func(t *testing.T) {
		cmd := NewReportCmd(viper.New())
		outBuf := bytes.NewBufferString("")
		cmd.SetOut(outBuf)
		errBuf := bytes.NewBufferString("")
		cmd.SetErr(errBuf)

		cmd.SetArgs([]string{
			"--set", fmt.Sprintf("%s=%s", profiles.VersionConfigName, "2.1"),
			report.MetadataCommandName,
			"test/report.yaml",
		})
		require.NoError(t, cmd.Execute())

		testReport := report.OutputReport{}
		require.NoError(t, json.Unmarshal([]byte(outBuf.String()), &testReport))

		require.True(t, compareMetadata(expectedMetadata, testReport.MetadataReport))

	})

}

func compareMetadata(expected *report.MetadataReport, result *report.MetadataReport) bool {
	outcome := true
	if strings.Compare(expected.ProfileVersion, result.ProfileVersion) != 0 {
		fmt.Println(fmt.Sprintf("profile version mistmatch %s : %s", expected.ProfileVersion, result.ProfileVersion))
		outcome = false
	}
	if expected.ProfileVendorType != result.ProfileVendorType {
		fmt.Println(fmt.Sprintf("profile vendortype mistmatch %s : %s", expected.ProfileVendorType, result.ProfileVendorType))
		outcome = false
	}
	if expected.ChartUri != result.ChartUri {
		fmt.Println(fmt.Sprintf("chart uri mistmatch %s : %s", expected.ChartUri, result.ChartUri))
		outcome = false
	}
	if expected.Chart.Name != result.Chart.Name {
		fmt.Println(fmt.Sprintf("chart name mistmatch %s : %s", expected.Chart.Name, result.Chart.Name))
		outcome = false
	}
	if expected.Chart.Version != result.Chart.Version {
		fmt.Println(fmt.Sprintf("chart version mistmatch %s : %s", expected.Chart.Version, result.Chart.Version))
		outcome = false
	}
	return outcome
}

func compareDigests(expected *report.DigestReport, result *report.DigestReport) bool {
	outcome := true
	if strings.Compare(expected.PackageDigest, result.PackageDigest) != 0 {
		fmt.Println(fmt.Sprintf("package digest mistmatch %s : %s", expected.PackageDigest, result.PackageDigest))
		outcome = false
	}
	if strings.Compare(expected.ChartDigest, result.ChartDigest) != 0 {
		fmt.Println(fmt.Sprintf("chart digest mistmatch %s : %s", expected.ChartDigest, result.ChartDigest))
		outcome = false
	}
	return outcome
}

func compareResults(expected *report.ResultsReport, result *report.ResultsReport) bool {
	outcome := true
	if strings.Compare(expected.Passed, result.Passed) != 0 {
		fmt.Println(fmt.Sprintf("results passed mistmatch %s : %s", expected.Passed, result.Passed))
		outcome = false
	}
	if strings.Compare(expected.Failed, result.Failed) != 0 {
		fmt.Println(fmt.Sprintf("results failed mistmatch %s : %s", expected.Failed, result.Failed))
		outcome = false
	}
	numMessages, err := strconv.Atoi(result.Failed)
	if err != nil {
		fmt.Println(fmt.Sprintf("results failed cannot be converted to int  %s : %v", result.Failed, err))
		outcome = false
	} else if len(result.Messages) != numMessages {
		fmt.Println(fmt.Sprintf("results number of fails and number of messages mismatch %d : %d", len(result.Messages), numMessages))
		outcome = false
	}
	return outcome
}

func compareAnnotations(expected []report.Annotation, result []report.Annotation) bool {
	outcome := true
	if len(expected) != len(result) {
		fmt.Println(fmt.Sprintf("num of annotation mismtatch %d : %d", len(expected), len(result)))
		outcome = false
	}
	for _, expectedAnnotation := range expected {
		found := false
		for _, resultAnnotation := range result {
			if strings.Compare(expectedAnnotation.Name, resultAnnotation.Name) == 0 {
				found = true
				if strings.Compare(expectedAnnotation.Value, resultAnnotation.Value) != 0 {
					fmt.Println(fmt.Sprintf("%s annotation mismtatch %s : %s", expectedAnnotation.Name, expectedAnnotation.Value, resultAnnotation.Value))
					outcome = false
				}
			}
		}
		if !found {
			fmt.Println(fmt.Sprintf("%s annotation not found in results", expectedAnnotation.Name))
			outcome = false
		}
	}
	return outcome
}
