package tool

import (
	"fmt"
	"os"
	"testing"

	"github.com/helm/chart-testing/v3/pkg/exec"
	"github.com/stretchr/testify/require"
)

func temporarilyToggleEnv(name string) (untoggleFn func(), err error) {
	originalValue := os.Getenv(name)
	err = os.Unsetenv(name)
	untoggleFn = func() {
		os.Setenv(name, originalValue)
	}
	return
}

func TestInstallWithValues(t *testing.T) {
	t.Run("helm not available in path", func(t *testing.T) {
		untogglePathEnv, err := temporarilyToggleEnv("PATH")
		require.NoError(t, err)
		defer untogglePathEnv()

		// arrange
		valuesFile := ""
		chrt := "non-existing-chart.tgz"
		release := "non-existing-release"
		namespace := "default"
		expectedErrorMessage := fmt.Sprintf(
			"executing helm with args \"install %s %s --namespace %s --wait\": "+
				"Error running process: exec: \"helm\": executable file not found in $PATH",
			release, chrt, namespace)
		h := NewHelm(exec.NewProcessExecutor(false), []string{})

		// act
		// for this test, none of the arguments matter si
		err = h.InstallWithValues(chrt, valuesFile, namespace, release)

		// assert
		require.Error(t, err)
		require.Equal(t, expectedErrorMessage, err.Error())
	})

	t.Run("helm failure should include content streamed to both Stderr and Stdout", func(t *testing.T) {
		// arrange
		chrt := "non-existing-chart.tgz"
		release := "non-existing-release"
		valuesFile := ""
		namespace := "default"
		expectedErrorMessage := fmt.Sprintf(
			"executing helm with args \"install %s %s --namespace %s --wait\": "+
				"Error running process: exit status 1\n---\nError: failed to download \"%s\" (hint: running `helm repo update` may help)",
			release, chrt, namespace, chrt)
		processExecutor := exec.NewProcessExecutor(false)
		extraArgs := []string{}
		h := NewHelm(processExecutor, extraArgs)

		// act
		err := h.InstallWithValues(chrt, valuesFile, namespace, release)

        // assert
		require.Error(t, err)
		require.Equal(t, expectedErrorMessage, err.Error())
	})
}
