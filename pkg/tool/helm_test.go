package tool

import (
	"testing"

	"github.com/helm/chart-testing/v3/pkg/exec"
	"github.com/stretchr/testify/require"
)

func TestInstallWithValues(t *testing.T) {
	t.Run("failure with empty Stdout capture should not include detail block", func(t *testing.T) {
		// arrange
		processExecutor := exec.NewProcessExecutor(false)
		extraArgs := []string{}
		h := NewHelm(processExecutor, extraArgs)

		// act
		chrt := "non-existing-chart.tgz"
		release := "non-existing-release"
		valuesFile := ""
		namespace := "default"
		err := h.InstallWithValues(chrt, valuesFile, namespace, release)

		// assert
		require.Error(t, err)

        // 
        require.Equal(
            t, "executing helm with args \"install non-existing-release non-existing-chart.tgz --namespace default --wait\": " +
                "Error running process: exit status 1",
            err.Error())
	})
}
