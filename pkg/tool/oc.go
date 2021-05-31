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

const osVersionKey = "serverVersion"

// Based on https://access.redhat.com/solutions/4870701
var kubeOpenShiftVersionMap map[string]string = map[string]string{
	"1.20": "4.7.0",
	"1.19": "4.6.0",
	"1.18": "4.5.0",
	"1.17": "4.4.0",
	"1.16": "4.3.0",
	"1.14": "4.2.0",
	"1.13": "4.1.0",
}

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
	// Relying on Kubernetes version can be replaced after fixing this issue:
	// https://bugzilla.redhat.com/show_bug.cgi?id=1850656
	kubeServerVersion := out[osVersionKey].(map[string]string)
	kubeVersion := kubeServerVersion["major"] + "." + kubeServerVersion["minor"]
	osVersion, ok := kubeOpenShiftVersionMap[kubeVersion]
	if !ok {
		return "", fmt.Errorf("Internal error: %q not found in Kubernetes-OpenShift version map", kubeVersion)
	}

	return osVersion, nil
}
