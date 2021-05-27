package tool

import (
	"fmt"
	"strings"

	"github.com/helm/chart-testing/v3/pkg/exec"
	"github.com/helm/chart-testing/v3/pkg/tool"
)

// Helm is an interface to the helm binary; it is a thin layer on top of the Helm abstraction offered by chart-testing
// to silence output being streamed to Stdout.
type Helm struct {
	tool.Helm
	exec.ProcessExecutor
	extraArgs []string
}

func NewHelm(exec exec.ProcessExecutor, extraArgs []string) Helm {
	return Helm{
		tool.NewHelm(exec, extraArgs),
		exec,
		extraArgs,
	}
}

func (h Helm) InstallWithValues(chart string, valuesFile string, namespace string, release string) error {
	var values []string
	if valuesFile != "" {
		values = []string{"--values", valuesFile}
	}

	helmInstallArgs := []string{"install", release, chart, "--namespace", namespace, "--wait"}
	helmInstallArgs = append(helmInstallArgs, values...)
	helmInstallArgs = append(helmInstallArgs, h.extraArgs...)

	if _, err := h.RunProcessAndCaptureOutput("helm", helmInstallArgs); err != nil {
		return fmt.Errorf("executing helm with args %q: %w", strings.Join(helmInstallArgs, " "), err)
	}

	return nil
}

func (h Helm) Test(namespace string, release string) error {
	_, err := h.RunProcessAndCaptureOutput("helm", "test", release, "--namespace", namespace, h.extraArgs)
	return err
}

func (h Helm) DeleteRelease(namespace string, release string) {
	_, _ = h.RunProcessAndCaptureOutput("helm", "uninstall", release, "--namespace", namespace, h.extraArgs)

}
