package tool

import "testing"

type fakeProcessExecutor struct{}

var output = `
serverVersion:
  major: "1"
  minor: "20"
`

func (p fakeProcessExecutor) RunProcessAndCaptureOutput(executable string, execArgs ...interface{}) (string, error) {
	return output, nil
}

func TestOcVersion(t *testing.T) {
	fp := fakeProcessExecutor{}
	oc := NewOc(fp)
	version, _ := oc.GetVersion()
	if version != "4.7.0" {
		t.Error("Version mismatch", version)
	}
}
