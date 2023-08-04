package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	helmchart "helm.sh/helm/v3/pkg/chart"

	"github.com/redhat-certification/chart-verifier/internal/chartverifier/profiles"
	"github.com/redhat-certification/chart-verifier/internal/chartverifier/utils"
	apireportsummary "github.com/redhat-certification/chart-verifier/pkg/chartverifier/reportsummary"
)

func TestReport(t *testing.T) {
	var (
		expectedResults = &apireportsummary.ResultsReport{
			Passed: "11",
			Failed: "1",
		}

		annotationPrefix = "charts.testing.io"

		expectedMetadata = &apireportsummary.MetadataReport{
			ProfileVersion:    "v1.1",
			ProfileVendorType: "redhat",
			ChartUri:          "internal/chartverifier/checks/chart-0.1.0-v3.valid.tgz",
			Chart:             &helmchart.Metadata{Name: "chart", Version: "0.1.0-v3.valid"},
			WebCatalogOnly:    false,
		}

		expectedDigests = &apireportsummary.DigestReport{
			PackageDigest: "4f29f2a95bf2b9a1c62fd215b079a01bdc5a38e9b4ff874d0fa21d0afca2e76d",
			ChartDigest:   "sha256:0c1c44def5c5de45212d90396062e18e0311b07789f477268fbf233c1783dbd0",
		}

		expectedAnnotations = []apireportsummary.Annotation{
			{
				Name:  fmt.Sprintf("%s/%s", apireportsummary.DefaultAnnotationsPrefix, apireportsummary.DigestsAnnotationName),
				Value: "sha256:0c1c44def5c5de45212d90396062e18e0311b07789f477268fbf233c1783dbd0",
			},
			{
				Name:  fmt.Sprintf("%s/%s", apireportsummary.DefaultAnnotationsPrefix, apireportsummary.TestedOCPVersionAnnotationName),
				Value: "4.7.8",
			},
			{
				Name:  fmt.Sprintf("%s/%s", apireportsummary.DefaultAnnotationsPrefix, apireportsummary.LastCertifiedTimestampAnnotationName),
				Value: "2021-07-06T10:28:01.09604-04:00",
			},
			{
				Name:  fmt.Sprintf("%s/%s", apireportsummary.DefaultAnnotationsPrefix, apireportsummary.SupportedOCPVersionsAnnotationName),
				Value: "4.7.8",
			},
		}
	)

	tests := []struct {
		name            string
		args            []string
		wantErr         bool
		wantResults     *apireportsummary.ResultsReport
		wantMetadata    *apireportsummary.MetadataReport
		wantAnnotations []apireportsummary.Annotation
		wantDigests     *apireportsummary.DigestReport
	}{
		{
			name:    "Should fail when no argument is given",
			args:    []string{},
			wantErr: true,
		},
		{
			name:    "Should fail when one argument is given",
			args:    []string{"test/report.yaml"},
			wantErr: true,
		},
		{
			name: "Should fail when bad subcommand is given",
			args: []string{
				"None",
				"test/report.yaml",
			},
			wantErr: true,
		},
		{
			name: "Should pass for subcommand annotations",
			args: []string{
				string(apireportsummary.AnnotationsSummary),
				"test/report.yaml",
			},
			wantAnnotations: expectedAnnotations,
		},
		{
			name: "Should pass for subcommand results",
			args: []string{
				string(apireportsummary.ResultsSummary),
				"test/report.yaml",
			},
			wantResults: expectedResults,
		},
		{
			name: "Should pass for subcommand metadata",
			args: []string{
				string(apireportsummary.MetadataSummary),
				"test/report.yaml",
			},
			wantMetadata: expectedMetadata,
		},
		{
			name: "Should pass for subcommand digests",
			args: []string{
				string(apireportsummary.DigestsSummary),
				"test/report.yaml",
			},
			wantDigests: expectedDigests,
		},
		{
			name: "Should pass for subcommand all",
			args: []string{
				string(apireportsummary.AllSummary),
				"test/report.yaml",
			},
			wantDigests:     expectedDigests,
			wantAnnotations: expectedAnnotations,
			wantMetadata:    expectedMetadata,
			wantResults:     expectedResults,
		},
		{
			name: "Should pass for annotation prefix",
			args: []string{
				"--set", fmt.Sprintf("%s=%s", apireportsummary.AnnotationsPrefixConfigName, annotationPrefix),
				string(apireportsummary.AnnotationsSummary),
				"test/report.yaml",
			},
			wantAnnotations: []apireportsummary.Annotation{
				{
					Name:  fmt.Sprintf("%s/%s", annotationPrefix, apireportsummary.DigestsAnnotationName),
					Value: "sha256:0c1c44def5c5de45212d90396062e18e0311b07789f477268fbf233c1783dbd0",
				},
				{
					Name:  fmt.Sprintf("%s/%s", annotationPrefix, apireportsummary.TestedOCPVersionAnnotationName),
					Value: "4.7.8",
				},
				{
					Name:  fmt.Sprintf("%s/%s", annotationPrefix, apireportsummary.LastCertifiedTimestampAnnotationName),
					Value: "2021-07-06T10:28:01.09604-04:00",
				},
				{
					Name:  fmt.Sprintf("%s/%s", annotationPrefix, apireportsummary.SupportedOCPVersionsAnnotationName),
					Value: "4.7.8",
				},
			},
		},
		{
			name: "Should pass for community profile",
			args: []string{
				"--set", fmt.Sprintf("%s=%s", profiles.VendorTypeConfigName, "community"),
				string(apireportsummary.AllSummary),
				"test/report.yaml",
			},
			wantMetadata: expectedMetadata,
			wantResults: &apireportsummary.ResultsReport{
				Passed: "1",
				Failed: "0",
			},
		},
		{
			name: "Should pass for invalid profile version",
			args: []string{
				"--set", fmt.Sprintf("%s=%s", profiles.VersionConfigName, "2.1"),
				string(apireportsummary.MetadataSummary),
				"test/report.yaml",
			},
			wantMetadata: expectedMetadata,
		},
		{
			name: "Should pass writing output to yaml file",
			args: []string{
				"-w",
				"-o",
				"yaml",
				string(apireportsummary.AnnotationsSummary),
				"test/report.yaml",
			},
			wantErr: false,
		},
		{
			name: "Should pass writing output to json file",
			args: []string{
				"-w",
				"-o",
				"json",
				string(apireportsummary.AnnotationsSummary),
				"test/report.yaml",
			},
			wantErr: false,
		},
		{
			name: "Should pass for skip digest check",
			args: []string{
				"-d",
				string(apireportsummary.MetadataSummary),
				"test/badDigestReport.yaml",
			},
			wantMetadata: expectedMetadata,
		},
		{
			name: "Should error with bad digest check",
			args: []string{
				string(apireportsummary.MetadataSummary),
				"test/badDigestReport.yaml",
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		outBuff := bytes.NewBufferString("")
		errBuff := bytes.NewBufferString("")
		var testReport apireportsummary.ReportSummary

		hasOutput := func() bool {
			return len(tc.wantAnnotations) > 0 ||
				tc.wantDigests != nil ||
				tc.wantMetadata != nil ||
				tc.wantResults != nil
		}()

		t.Run(tc.name, func(t *testing.T) {
			cmd := NewReportCmd(viper.New())
			cmd.SetOut(outBuff)
			if hasOutput {
				utils.CmdStdout = outBuff
			}
			cmd.SetErr(errBuff)

			if len(tc.args) > 0 {
				cmd.SetArgs(tc.args)
			}

			if tc.wantErr {
				require.Error(t, cmd.Execute())
			} else {
				require.NoError(t, cmd.Execute())
			}

			if hasOutput {
				require.NoError(t, json.Unmarshal(outBuff.Bytes(), &testReport))
			}

			if len(tc.wantAnnotations) > 0 {
				require.True(t, compareAnnotations(tc.wantAnnotations, testReport.AnnotationsReport))
			}

			if tc.wantResults != nil {
				require.True(t, compareResults(tc.wantResults, testReport.ResultsReport))
			}

			if tc.wantDigests != nil {
				require.True(t, compareDigests(tc.wantDigests, testReport.DigestsReport))
			}

			if tc.wantMetadata != nil {
				require.True(t, compareMetadata(tc.wantMetadata, testReport.MetadataReport))
			}
		})
	}
}

func compareMetadata(expected *apireportsummary.MetadataReport, result *apireportsummary.MetadataReport) bool {
	outcome := true
	if strings.Compare(expected.ProfileVersion, result.ProfileVersion) != 0 {
		fmt.Printf("profile version mismatch %s : %s\n", expected.ProfileVersion, result.ProfileVersion)
		outcome = false
	}
	if expected.ProfileVendorType != result.ProfileVendorType {
		fmt.Printf("profile vendortype mismatch %s : %s\n", expected.ProfileVendorType, result.ProfileVendorType)
		outcome = false
	}
	if expected.ChartUri != result.ChartUri {
		fmt.Printf("chart uri mismatch %s : %s\n", expected.ChartUri, result.ChartUri)
		outcome = false
	}
	if expected.Chart.Name != result.Chart.Name {
		fmt.Printf("chart name mismatch %s : %s\n", expected.Chart.Name, result.Chart.Name)
		outcome = false
	}
	if expected.Chart.Version != result.Chart.Version {
		fmt.Printf("chart version mismatch %s : %s\n", expected.Chart.Version, result.Chart.Version)
		outcome = false
	}
	if expected.WebCatalogOnly != result.WebCatalogOnly {
		fmt.Printf("web catalog only mismatch %t : %t\n", expected.WebCatalogOnly, result.WebCatalogOnly)
		outcome = false
	}

	return outcome
}

func compareDigests(expected *apireportsummary.DigestReport, result *apireportsummary.DigestReport) bool {
	outcome := true
	if strings.Compare(expected.PackageDigest, result.PackageDigest) != 0 {
		fmt.Printf("package digest mismatch %s : %s\n", expected.PackageDigest, result.PackageDigest)
		outcome = false
	}
	if strings.Compare(expected.ChartDigest, result.ChartDigest) != 0 {
		fmt.Printf("chart digest mismatch %s : %s\n", expected.ChartDigest, result.ChartDigest)
		outcome = false
	}
	return outcome
}

func compareResults(expected *apireportsummary.ResultsReport, result *apireportsummary.ResultsReport) bool {
	outcome := true
	if strings.Compare(expected.Passed, result.Passed) != 0 {
		fmt.Printf("results passed mismatch; want %s but got %s\n", expected.Passed, result.Passed)
		outcome = false
	}
	if strings.Compare(expected.Failed, result.Failed) != 0 {
		fmt.Printf("results failed mismatch %s : %s\n", expected.Failed, result.Failed)
		outcome = false
	}
	numMessages, err := strconv.Atoi(result.Failed)
	if err != nil {
		fmt.Printf("results failed cannot be converted to int  %s : %v\n", result.Failed, err)
		outcome = false
	} else if len(result.Messages) != numMessages {
		fmt.Printf("results number of fails and number of messages mismatch %d : %d\n", len(result.Messages), numMessages)
		outcome = false
	}
	return outcome
}

func compareAnnotations(expected []apireportsummary.Annotation, result []apireportsummary.Annotation) bool {
	outcome := true
	if len(expected) != len(result) {
		fmt.Printf("num of annotation mismatch %d : %d\n", len(expected), len(result))
		outcome = false
	}

	for _, expectedAnnotation := range expected {
		found := false
		for _, resultAnnotation := range result {
			if strings.Compare(expectedAnnotation.Name, resultAnnotation.Name) == 0 {
				found = true
				if strings.Compare(expectedAnnotation.Value, resultAnnotation.Value) != 0 {
					fmt.Printf("%s annotation mismatch %s : %s\n", expectedAnnotation.Name, expectedAnnotation.Value, resultAnnotation.Value)
					outcome = false
				}
			}
		}
		if !found {
			fmt.Printf("%s annotation not found in results\n", expectedAnnotation.Name)
			outcome = false
		}
	}
	return outcome
}
