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

func (h Helm) RunProcessAndCaptureOutput(executable string, execArgs ...interface{}) (string, error) {
	return h.RunProcessInDirAndCaptureOutput("", executable, execArgs...)
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

// RunProcessInDirAndCaptureOutput overrides exec.ProcessExecutor's and inject the command line and any streamed content
// to either Stdout or Stderr into the returned error, if any.
func (h Helm) RunProcessInDirAndCaptureOutput(
    workingDirectory string,
    executable string,
    execArgs ...interface{},
) (string, error) {
	cmd, err := h.CreateProcess(executable, execArgs...)
	if err != nil {
		return "", err
	}

	cmd.Dir = workingDirectory
	bytes, err := cmd.CombinedOutput()
	capturedOutput := strings.TrimSpace(string(bytes))

	execArgsCopy := toStringArray(execArgs)
	execArgsStr := strings.Join(execArgsCopy, " ")

	if err != nil {
		if len(capturedOutput) == 0 {
			return "", fmt.Errorf(
                "Error running process: executing %s with args %q: %w",
                executable, execArgsStr, err)
		}
		return capturedOutput, fmt.Errorf(
            "Error running process: executing %s with args %q: %w\n---\n%s",
            executable, execArgsStr, err, capturedOutput)
	}
	return capturedOutput, nil
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
