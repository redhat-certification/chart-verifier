package tool

import (
	"fmt"
	"testing"
)

type TestData struct {
	getVersionOut string
	OCVersion     string
}

type fakeProcessExecutor struct {
	Output string
}

var output120 = `
serverVersion:
  major: "1"
  minor: "20"
`
var output121 = `
serverVersion:
  major: "1"
  minor: "21"
`
var output122 = `
serverVersion:
  major: "1"
  minor: "22"
`

var testsData []TestData

func (p fakeProcessExecutor) RunProcessAndCaptureOutput(executable string, execArgs ...interface{}) (string, error) {
	return p.Output, nil
}

func TestOCVersions(t *testing.T) {

	testsData = append(testsData, TestData{getVersionOut: output120, OCVersion: "4.7.0"})
	testsData = append(testsData, TestData{getVersionOut: output121, OCVersion: "4.8.0"})
	testsData = append(testsData, TestData{getVersionOut: output122, OCVersion: "4.9.0"})

	for _, testdata := range testsData {
		t.Logf("Check for OC %s", testdata.OCVersion)
		checkOcVersion(testdata.getVersionOut, testdata.OCVersion, t)
	}
}

func checkOcVersion(kubeVersion string, expectedOCversion string, t *testing.T) {
	fp := fakeProcessExecutor{Output: kubeVersion}
	oc := NewOc(fp)
	version, _ := oc.GetVersion()
	if version != expectedOCversion {
		t.Error(fmt.Sprintf("Version mismatch expected: %s, got: %s", expectedOCversion, version))
	}
}
