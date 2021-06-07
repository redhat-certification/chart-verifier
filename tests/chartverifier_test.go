package hack

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/redhat-certification/chart-verifier/pkg/tool"
	"github.com/stretchr/testify/require"
)

// absPathFromSourceFileLocation returns the absolute path of a file or directory under the current source file's
// location.
//
// TODO: refactor this into a testutil package.
func absPathFromSourceFileLocation(name string) (string, error) {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		panic("couldn't get current path")
	}
	filename, err := filepath.Abs(filename)
	if err != nil {
		return "", fmt.Errorf("retrieving current source file's location: %w", err)
	}
	dirname := path.Dir(filename)
	return filepath.Join(dirname, name), nil
}

func TestChartVerifierSh(t *testing.T) {
	t.Run("Should succeed when the chart exists and is valid for a single check", func(t *testing.T) {
		if os.Getenv("CHART_VERIFIER_ENABLE_CLUSTER_TESTING") == "" {
			t.Skip("CHART_VERIFIER_ENABLE_CLUSTER_TESTING not set, skipping in cluster tests")
		}

		shPath, err := absPathFromSourceFileLocation("../hack/chart-verifier.sh")
		require.NoError(t, err)

		pkgPath, err := absPathFromSourceFileLocation("../pkg/chartverifier/checks/psql-service-0.1.7")
		require.NoError(t, err)

		args := []interface{}{
			"-V", "4.8",
			"--chart-set", "k8Project=default",
			pkgPath,
		}

		pe := tool.NewProcessExecutor(false)
		cmdOutput, err := pe.RunProcessInDirAndCaptureOutput(".", shPath, args...)

		require.NoError(t, err)
		require.NotEmpty(t, cmdOutput)

		expected := "- check: is-helm-v3\n" +
			"    type: Mandatory\n" +
			"    outcome: PASS\n"

		require.Contains(t, cmdOutput, expected)
	})
}
