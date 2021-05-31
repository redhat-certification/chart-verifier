package tool

import (
	"testing"

	"github.com/helm/chart-testing/v3/pkg/exec"
	"github.com/stretchr/testify/require"
)

func TestInstallWithValues(t *testing.T) {
    // arrange
    processExecutor := exec.NewProcessExecutor(false)
    extraArgs := []string{}
	h := NewHelm(processExecutor, extraArgs)

    // act
    chrt := ""
    release := ""
    valuesFile := ""
    namespace := "default"
    err := h.InstallWithValues(chrt, valuesFile, namespace, release)

    // assert
    require.Error(t, err)
    
}
