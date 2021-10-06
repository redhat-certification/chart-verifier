package tool

import (
	"fmt"
	"os"
	"strings"
	"testing"

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
			"Error running process: executing helm with args \"install %s %s --namespace %s --wait\": "+
				"exec: \"helm\": executable file not found in $PATH",
			release, chrt, namespace)
		h := NewHelm(NewProcessExecutor(false), []string{})

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
		expectedErrorMessage := "Error running process"
		processExecutor := NewProcessExecutor(false)
		extraArgs := []string{}
		h := NewHelm(processExecutor, extraArgs)

		// act
		err := h.InstallWithValues(chrt, valuesFile, namespace, release)

		// assert
		require.Error(t, err)
		require.True(t, strings.Contains(err.Error(), expectedErrorMessage))
	})
}
