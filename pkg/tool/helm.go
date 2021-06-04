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
    argsCopy := make([]string, len(args))
	for i, a := range args {
		argsCopy[i] = fmt.Sprint(a)
	}
	return argsCopy
}

func toInterfaceArray(args []string) []interface{} {
    argsCopy := make([]interface{}, len(args))
	for i, a := range args {
		argsCopy[i] = a
	}
	return argsCopy
}

// InstallWithValues overrides chart-testing's tool.Helm method to execute the modified RunProcessAndCaptureOutput
// method.
func (h Helm) InstallWithValues(chart string, valuesFile string, namespace string, release string) error {
	var values []interface{}
	if valuesFile != "" {
		values = []interface{}{"--values", valuesFile}
	}

	helmArgs := []interface{}{"install", release, chart, "--namespace", namespace, "--wait"}
	helmArgs = append(helmArgs, values...)
	helmArgs = append(helmArgs,  toInterfaceArray(h.extraArgs)...)

	_, err := h.RunProcessAndCaptureOutput("helm", helmArgs...)
	return err
}

func (h Helm) Test(namespace string, release string) error {
	_, err := h.RunProcessAndCaptureOutput("helm", "test", release, "--namespace", namespace, h.extraArgs)
	return err
}

func (h Helm) DeleteRelease(namespace string, release string) {
	_, _ = h.RunProcessAndCaptureOutput("helm", "uninstall", release, "--namespace", namespace, h.extraArgs)

}
