package hack

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChartVerifierSh(t *testing.T) {
	t.Run("Should succeed when the chart exists and is valid for a single check", func(t *testing.T) {
		if os.Getenv("CHART_VERIFIER_ENABLE_CLUSTER_TESTING") == "" {
			t.Skip("CHART_VERIFIER_ENABLE_CLUSTER_TESTING not set, skipping in cluster tests")
		}
		t.Skip()
		args := []string{
			"-V", "4.9",
			"../pkg/chartverifier/checks/chart-0.1.0-v3.valid.tgz",
		}

		cmdOutput, err := exec.Command("./chart-verifier.sh", args...).Output()

		require.NoError(t, err)
		require.NotEmpty(t, cmdOutput)

		expected := "results:\n" +
			"  - check: is-helm-v3\n" +
			"    type: Mandatory\n" +
			"    outcome: PASS\n" +
			"    reason: API version is V2, used in Helm 3\n"

		require.Contains(t, cmdOutput, expected)
	})
}
