package tool

import (
	"fmt"
	"strings"

	"github.com/helm/chart-testing/v3/pkg/exec"
	"github.com/helm/chart-testing/v3/pkg/tool"
)

// Kubectl is an interface to the helm binary; it is a thin layer on top of the Kubectl abstraction offered by
// chart-testing to silence output being streamed to Stdout.
type Kubectl struct {
	tool.Kubectl
	exec.ProcessExecutor
}

func NewKubectl(exec exec.ProcessExecutor) Kubectl {
	return Kubectl{
		tool.NewKubectl(exec),
		exec,
	}
}

func (k Kubectl) WaitForDeployments(namespace string, selector string) error {
	output, err := k.RunProcessAndCaptureOutput(
		"kubectl", "get", "deployments", "--namespace", namespace, "--selector", selector, "--output", "jsonpath={.items[*].metadata.name}")
	if err != nil {
		return err
	}

	deployments := strings.Fields(output)
	for _, deployment := range deployments {
		deployment = strings.Trim(deployment, "'")
		_, err := k.RunProcessAndCaptureOutput("kubectl", "rollout", "status", "deployment", deployment, "--namespace", namespace)
		if err != nil {
			return err
		}

		// 'kubectl rollout status' does not return a non-zero exit code when rollouts fail.
		// We, thus, need to double-check here.
		//
		// Just after rollout, pods from the previous deployment revision may still be in a
		// terminating state.
		unavailable, err := k.RunProcessAndCaptureOutput("kubectl", "get", "deployment", deployment, "--namespace", namespace, "--output",
			`jsonpath={.status.unavailableReplicas}`)
		if err != nil {
			return err
		}
		if unavailable != "" && unavailable != "0" {
			return fmt.Errorf("%s replicas unavailable", unavailable)
		}
	}

	return nil
}
