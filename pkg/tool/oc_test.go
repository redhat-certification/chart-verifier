package tool

import "testing"

type fakeProcessExecutor struct {
	Output string
}

var output47 = `
serverVersion:
  major: "1"
  minor: "20"
`

func (p fakeProcessExecutor) RunProcessAndCaptureOutput(executable string, execArgs ...interface{}) (string, error) {
	return p.Output, nil
}

func TestOcVersion47(t *testing.T) {
	fp := fakeProcessExecutor{Output: output47}
	oc := NewOc(fp)
	version, _ := oc.GetVersion()
	if version != "4.7.0" {
		t.Error("Version mismatch", version)
	}
}

var output48 = `
serverVersion:
  major: "1"
  minor: "21"
`

func TestOcVersion48(t *testing.T) {
	fp := fakeProcessExecutor{Output: output48}
	oc := NewOc(fp)
	version, _ := oc.GetVersion()
	if version != "4.8.0" {
		t.Error("Version mismatch", version)
	}
}
