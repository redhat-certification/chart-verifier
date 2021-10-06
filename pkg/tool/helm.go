package tool

import (
	"fmt"

	"github.com/helm/chart-testing/v3/pkg/tool"
)

// Helm is an interface to the helm binary; it is a thin layer on top of the Helm abstraction offered by chart-testing
// to silence output being streamed to Stdout.
type Helm struct {
	tool.Helm
	ProcessExecutor
	extraArgs []string
}

func NewHelm(exec ProcessExecutor, extraArgs []string) Helm {
	return Helm{
		tool.NewHelm(exec.ProcessExecutor, extraArgs),
		exec,
		extraArgs,
	}
}

func toStringArray(args []interface{}) []string {
	copy := make([]string, len(args))
	for i, a := range args {
		copy[i] = fmt.Sprint(a)
	}
	return copy
}

func toInterfaceArray(args []string) []interface{} {
	copy := make([]interface{}, len(args))
	for i, a := range args {
		copy[i] = a
	}
	return copy
}

// InstallWithValues overrides chart-testing's tool.Helm method to execute the modified RunProcessAndCaptureOutput
// method.
func (h Helm) InstallWithValues(chart string, valuesFile string, namespace string, release string) error {
	var values []interface{}
	if valuesFile != "" {
		values = []interface{}{"--values", valuesFile}
	}

	LogInfo(fmt.Sprintf("Execute helm install. namespace: %s, release: %s chart: %s", namespace, release, chart))
	helmArgs := []interface{}{"install", release, chart, "--namespace", namespace, "--wait"}
	helmArgs = append(helmArgs, values...)
	helmArgs = append(helmArgs, toInterfaceArray(h.extraArgs)...)

	_, err := h.RunProcessAndCaptureOutput("helm", helmArgs...)
	if err != nil {
		LogError(fmt.Sprintf("Execute helm install. error %v", err))
	} else {
		LogInfo("Helm install complete")
	}

	return err
}

func (h Helm) Test(namespace string, release string) error {
	LogInfo(fmt.Sprintf("Execute helm test. namespace: %s, release: %s, extraArgd: %v", namespace, release, h.extraArgs))
	_, err := h.RunProcessAndCaptureOutput("helm", "test", release, "--namespace", namespace, h.extraArgs)
	if err != nil {
		LogError(fmt.Sprintf("Execute helm test. error %v", err))
	} else {
		LogInfo("Helm test complete")
	}
	return err
}

func (h Helm) DeleteRelease(namespace string, release string) {
	LogInfo(fmt.Sprintf("Execute helm uninstall. namespace: %s, release: %s", namespace, release))
	_, err := h.RunProcessAndCaptureOutput("helm", "uninstall", release, "--namespace", namespace, h.extraArgs)
	if err != nil {
		LogError(fmt.Sprintf("Error from helm uninstall : %v", err))
	} else {
		LogInfo("Delete release complete")
	}
}
