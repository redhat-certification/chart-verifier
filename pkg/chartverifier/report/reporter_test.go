package report

import (
	"encoding/json"
	"fmt"
	"github.com/redhat-certification/chart-verifier/pkg/chartverifier/profiles"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"strconv"
	"strings"
	"testing"
)

func TestReports(t *testing.T) {

	testRedHatMetaDataReport := &MetadataReport{}
	testRedHatMetaDataReport.ProfileVersion = "v1.0"
	testRedHatMetaDataReport.ProfileVendorType = "redhat"

	testPartnerMetaDataReport := &MetadataReport{}
	testPartnerMetaDataReport.ProfileVersion = "v1.0"
	testPartnerMetaDataReport.ProfileVendorType = "partner"

	var testAnnotationsReport []Annotation
	testAnnotationsReport = append(testAnnotationsReport, Annotation{Name: fmt.Sprintf("charts.openshift.io/%s", DigestsAnnotationName), Value: "sha256:0c1c44def5c5de45212d90396062e18e0311b07789f477268fbf233c1783dbd0"})
	testAnnotationsReport = append(testAnnotationsReport, Annotation{Name: fmt.Sprintf("charts.openshift.io/%s", LastCertifiedTimestampAnnotationName), Value: "2021-07-02T08:09:56.881793-04:00"})
	testAnnotationsReport = append(testAnnotationsReport, Annotation{Name: fmt.Sprintf("charts.openshift.io/%s", CertifiedOCPVersionAnnotationName), Value: "4.7.8"})

	testDigestReport := &DigestReport{}
	testDigestReport.PackageDigest = "4f29f2a95bf2b9a1c62fd215b079a01bdc5a38e9b4ff874d0fa21d0afca2e76d"
	testDigestReport.ChartDigest = "sha256:0c1c44def5c5de45212d90396062e18e0311b07789f477268fbf233c1783dbd0"

	testResultsReport := &ResultsReport{}
	testResultsReport.Passed = "11"
	testResultsReport.Failed = "0"

	testAllReport := &OutputReport{}
	testAllReport.ResultsReport = testResultsReport
	testAllReport.MetadataReport = testRedHatMetaDataReport
	testAllReport.DigestsReport = testDigestReport
	testAllReport.AnnotationsReport = testAnnotationsReport

	type testInfo struct {
		description       string
		path              string
		annotationsPrefix string
		setVendorType     string
		expectedReport    *OutputReport
	}

	var tests []testInfo

	allgoodTestInfo := testInfo{}
	allgoodTestInfo.path = "testreports/reportallgood.yaml"
	allgoodTestInfo.description = fmt.Sprintf("test all good report %s", allgoodTestInfo.path)
	allgoodTestInfo.expectedReport = testAllReport
	tests = append(tests, allgoodTestInfo)

	missingMandatoryTestInfo := testInfo{}
	missingMandatoryTestInfo.path = "testreports/reportmissingmandatory.yaml"
	missingMandatoryTestInfo.description = fmt.Sprintf("test missing mandatory report %s", missingMandatoryTestInfo.path)
	missingMandatoryTestInfo.expectedReport = &OutputReport{}
	missingMandatoryTestInfo.expectedReport.ResultsReport = &ResultsReport{}
	missingMandatoryTestInfo.expectedReport.ResultsReport.Passed = "9"
	missingMandatoryTestInfo.expectedReport.ResultsReport.Failed = "2"
	missingMandatoryTestInfo.expectedReport.MetadataReport = testPartnerMetaDataReport
	tests = append(tests, missingMandatoryTestInfo)

	withFailureTestInfo := testInfo{}
	withFailureTestInfo.path = "testreports/reportwithfailure.yaml"
	withFailureTestInfo.description = fmt.Sprintf("test missing failures report %s", missingMandatoryTestInfo.path)
	withFailureTestInfo.expectedReport = &OutputReport{}
	missingMandatoryTestInfo.expectedReport.ResultsReport = &ResultsReport{}
	missingMandatoryTestInfo.expectedReport.ResultsReport.Passed = "9"
	missingMandatoryTestInfo.expectedReport.ResultsReport.Failed = "2"
	withFailureTestInfo.expectedReport.MetadataReport = testRedHatMetaDataReport
	tests = append(tests, withFailureTestInfo)

	allsortsTestInfo := testInfo{}
	allsortsTestInfo.path = "testreports/reportallsorts.yaml"
	allsortsTestInfo.description = fmt.Sprintf("test allsorts report %s", missingMandatoryTestInfo.path)
	allsortsTestInfo.expectedReport = &OutputReport{}
	allsortsTestInfo.expectedReport.ResultsReport = &ResultsReport{}
	allsortsTestInfo.expectedReport.ResultsReport.Passed = "8"
	allsortsTestInfo.expectedReport.ResultsReport.Failed = "3"
	allsortsTestInfo.expectedReport.DigestsReport = &DigestReport{}
	allsortsTestInfo.expectedReport.DigestsReport.ChartDigest = "sha256:0c1c44def5c5de45212d90396062e18e0311b07789f477268fbf233c1783dbd0"
	allsortsTestInfo.expectedReport.MetadataReport = testPartnerMetaDataReport
	tests = append(tests, allsortsTestInfo)

	setBehaviorTestInfo := testInfo{}
	setBehaviorTestInfo.path = "testreports/reportmissingmandatory.yaml"
	setBehaviorTestInfo.description = fmt.Sprintf("test set behvaior missing mandatory report %s", missingMandatoryTestInfo.path)
	setBehaviorTestInfo.setVendorType = "community"
	setBehaviorTestInfo.expectedReport = &OutputReport{}
	setBehaviorTestInfo.expectedReport.ResultsReport = &ResultsReport{}
	setBehaviorTestInfo.expectedReport.ResultsReport.Passed = "1"
	setBehaviorTestInfo.expectedReport.ResultsReport.Failed = "0"
	setBehaviorTestInfo.annotationsPrefix = "test.report.command.io"
	var setBehaviorAnnotationsReport []Annotation
	setBehaviorAnnotationsReport = append(setBehaviorAnnotationsReport, Annotation{Name: fmt.Sprintf("%s/%s", setBehaviorTestInfo.annotationsPrefix, DigestsAnnotationName), Value: "sha256:0c1c44def5c5de45212d90396062e18e0311b07789f477268fbf233c1783dbd0"})
	setBehaviorAnnotationsReport = append(setBehaviorAnnotationsReport, Annotation{Name: fmt.Sprintf("%s/%s", setBehaviorTestInfo.annotationsPrefix, LastCertifiedTimestampAnnotationName), Value: "2021-07-02T08:09:56.881793-04:00"})
	setBehaviorAnnotationsReport = append(setBehaviorAnnotationsReport, Annotation{Name: fmt.Sprintf("%s/%s", setBehaviorTestInfo.annotationsPrefix, CertifiedOCPVersionAnnotationName), Value: "4.7.8"})
	setBehaviorTestInfo.expectedReport.AnnotationsReport = setBehaviorAnnotationsReport
	setBehaviorTestInfo.expectedReport.MetadataReport = testPartnerMetaDataReport
	tests = append(tests, setBehaviorTestInfo)

	for _, test := range tests {

		t.Run("Report test : "+test.description, func(t *testing.T) {
			options := &ReportOptions{}

			options.URI = test.path
			options.ViperConfig = viper.New()
			if len(test.annotationsPrefix) > 0 {
				options.ViperConfig.Set(AnnotationsPrefixConfigName, test.annotationsPrefix)
			}
			if len(test.setVendorType) > 0 {
				options.ViperConfig.Set(profiles.VendorTypeConfigName, test.setVendorType)
			}

			all, err := ReportCommandRegistry().Get(AllCommandsName)(options)
			assert.NoError(t, err, "error getting All")
			if err == nil {
				if test.expectedReport.MetadataReport != nil {
					assert.True(t, CompareMetadata(test.expectedReport.MetadataReport, all.MetadataReport), "all report: Metadata does not match")
				}
				if test.expectedReport.DigestsReport != nil {
					assert.True(t, CompareDigests(test.expectedReport.DigestsReport, all.DigestsReport), "all report: Digests do not match")
				}
				if test.expectedReport.ResultsReport != nil {
					assert.True(t, CompareResults(test.expectedReport.ResultsReport, all.ResultsReport), "all report: Results do not match")
				}
				if len(test.expectedReport.AnnotationsReport) > 0 {
					assert.True(t, CompareAnnotations(test.expectedReport.AnnotationsReport, all.AnnotationsReport), "all report: Annotations do not match")
				}
				_, err = json.Marshal(all)
				assert.NoError(t, err, "All report is not valid json")
			}

			if test.expectedReport.MetadataReport != nil {
				metadata, err := ReportCommandRegistry().Get(MetadataCommandName)(options)
				assert.NoError(t, err, "error getting Metadata")
				if err == nil {
					outcome := CompareMetadata(test.expectedReport.MetadataReport, metadata.MetadataReport)
					assert.True(t, outcome, "Metadata report does not match")
					reportjson, err := json.Marshal(metadata.MetadataReport)
					assert.NoError(t, err, "Metadata report is not valid json")
					if err == nil && !outcome {
						fmt.Println(fmt.Sprintf("MetaDataReport :\n%s", reportjson))
					}
				}
			}

			if len(test.expectedReport.AnnotationsReport) > 0 {
				annotations, err := ReportCommandRegistry().Get(AnnotationsCommandName)(options)
				assert.NoError(t, err, "error getting Annotations")
				if err == nil {
					outcome := CompareAnnotations(test.expectedReport.AnnotationsReport, annotations.AnnotationsReport)
					assert.True(t, outcome, "Annotations report does not match")
					reportjson, err := json.Marshal(annotations.AnnotationsReport)
					assert.NoError(t, err, "Annotations report is not valid json")
					if err == nil && !outcome {
						fmt.Println(fmt.Sprintf("AnnotationsReport :\n%s", reportjson))
					}
				}
			}

			if test.expectedReport.ResultsReport != nil {
				results, err := ReportCommandRegistry().Get(ResultsCommandName)(options)
				assert.NoError(t, err, "error getting Results")
				if err == nil {
					outcome := CompareResults(test.expectedReport.ResultsReport, all.ResultsReport)
					assert.True(t, outcome, "Results report does not match")
					reportjson, err := json.Marshal(results.ResultsReport)
					assert.NoError(t, err, "Results report is not valid json")
					if err == nil && !outcome {
						fmt.Println(fmt.Sprintf("AResultsReport :\n%s", reportjson))
					}
				}
			}

			if test.expectedReport.DigestsReport != nil {
				digests, err := ReportCommandRegistry().Get(DigestsCommandName)(options)
				assert.NoError(t, err, "error getting Digests")
				if err == nil {
					outcome := CompareDigests(test.expectedReport.DigestsReport, all.DigestsReport)
					assert.True(t, outcome, "Digests report does not match")
					reportjson, err := json.Marshal(digests.DigestsReport)
					assert.NoError(t, err, "Digests report is not valid json")
					if err == nil && !outcome {
						fmt.Println(fmt.Sprintf("DigestsReport :\n%s", reportjson))
					}
				}
			}
		})
	}
}

func CompareMetadata(expected *MetadataReport, result *MetadataReport) bool {
	outcome := true
	if strings.Compare(expected.ProfileVersion, result.ProfileVersion) != 0 {
		fmt.Println(fmt.Sprintf("profile version mistmatch %s : %s", expected.ProfileVersion, result.ProfileVersion))
		outcome = false
	}
	if expected.ProfileVendorType != result.ProfileVendorType {
		fmt.Println(fmt.Sprintf("profile vendortype mistmatch %s : %s", expected.ProfileVendorType, result.ProfileVendorType))
		outcome = false
	}
	return outcome
}

func CompareDigests(expected *DigestReport, result *DigestReport) bool {
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

func CompareResults(expected *ResultsReport, result *ResultsReport) bool {
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

func CompareAnnotations(expected []Annotation, result []Annotation) bool {
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
