package tool

import (
	"fmt"

	"github.com/helm/chart-testing/v3/pkg/exec"
	"gopkg.in/yaml.v3"
)

type Oc struct {
	exec exec.ProcessExecutor
}

func NewOc(exec exec.ProcessExecutor) Oc {
	return Oc{
		exec: exec,
	}
}

const osVersionKey = "openshiftVersion"

func (o Oc) GetVersion() (string, error) {
	rawOutput, err := o.exec.RunProcessAndCaptureOutput("oc", "version", "-o", "yaml")
	if err != nil {
		return "", err
	}
	out := map[string]interface{}{}
	err = yaml.Unmarshal([]byte(rawOutput), &out)
	if err != nil {
		return "", err
	}

	version := out[osVersionKey]
	if version == nil {
		return "", fmt.Errorf("%q not found in 'oc version' output", osVersionKey)
	}

	v, ok := version.(string)
	if !ok {
		return "", fmt.Errorf("%q does not contain a string: %v", osVersionKey, version)
	}

	return v, nil
}
