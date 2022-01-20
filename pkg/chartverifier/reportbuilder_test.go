package chartverifier

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

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

func TestFilePackageDigest(t *testing.T) {

	charts := []string{"checks/chart-0.1.0-v3.with-csi.tgz",
		"checks/chart-0.1.0-v3.valid.tgz",
		"checks/chart-0.1.0-v2.invalid.tgz",
		"checks/chart-0.1.0-v3.without-readme.tgz",
		"checks/chart-0.1.0-v3.with-crd.tgz"}

	for _, chart := range charts {
		cmd := exec.Command("sha256sum", chart)
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		assert.NoError(t, err, "error running sha256sum command")
		commandResponse := strings.Split(out.String(), " ")
		assert.Equal(t, commandResponse[0], GetPackageDigest(chart), fmt.Sprintf("%s digests did not match as expected", chart))
	}

}

func TestUrlPackageDigest(t *testing.T) {

	charts := make(map[string]string)

	charts["https://github.com/openshift-helm-charts/charts/releases/download/hashicorp-vault-0.13.0/hashicorp-vault-0.13.0.tgz"] = "97e274069d9d3d028903610a3f9fca892b2620f0a334de6215ec5f962328586f"
	charts["https://github.com/openshift-helm-charts/charts/releases/download/hashicorp-vault-0.12.0/hashicorp-vault-0.12.0.tgz"] = "b07be2a554ecbe6a6dd48ea763ed568de317d17cf1a19fb11ddb562983286555"
	charts["https://github.com/IBM/charts/blob/master/repo/ibm-helm/ibm-object-storage-plugin-2.1.2.tgz?raw=true"] = "06efa1e26f8a7ba93a6e6136650b0624af2558cc44a4588198fca322f9219e32"
	charts["checks/chart-0.1.0-v3.valid.tgz?raw=true"] = "4d6a38386eb8f3bbcdb4c1a4a6c3ccb7e8f38317e2a7924c0666087ff9b29c39"

	for chart, sha := range charts {

		assert.Equal(t, sha, GetPackageDigest(chart), fmt.Sprintf("%s digests did not match as expected", chart))

	}

}
