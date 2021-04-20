package chartverifier

import (
	"testing"

	"github.com/redhat-certification/chart-verifier/pkg/chartverifier/checks"
	"github.com/stretchr/testify/require"
)

func TestSha(t *testing.T) {

	charts := []string{"checks/chart-0.1.0-v3.with-csi.tgz",
		"checks/chart-0.1.0-v3.valid.tgz",
		"checks/chart-0.1.0-v2.invalid.tgz",
		"checks/chart-0.1.0-v3.without-readme.tgz",
		"checks/chart-0.1.0-v3.with-crd.tgz"}

	shas := make(map[string]string)

	// get shas for each file
	for _, chart := range charts {
		t.Run("Sha generation : "+chart, func(t *testing.T) {
			helmChart, _, err := checks.LoadChartFromURI(chart)
			require.NoError(t, err)
			sha := GenerateSha(helmChart.Raw)
			require.NotNil(t, sha)
			shas[chart] = sha
		})
	}

	for i := 0; i < 5; i++ {
		// get sha's again and make sure they match
		for _, chart := range charts {

			t.Run("Sha must not change : "+chart, func(t *testing.T) {
				helmChart, _, err := checks.LoadChartFromURI(chart)
				require.NoError(t, err)
				sha := GenerateSha(helmChart.Raw)
				require.NotNil(t, sha)
				require.Equal(t, shas[chart], sha)
			})
		}
	}

}
