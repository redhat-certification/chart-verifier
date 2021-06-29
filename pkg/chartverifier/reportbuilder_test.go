package chartverifier

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os/exec"
	"strings"
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

	charts := []string{"https://github.com/redhat-certification/chart-verifier/blob/main/pkg/chartverifier/checks/chart-0.1.0-v3.valid.tgz?raw=true",
		"https://github.com/openshift-helm-charts/charts/releases/download/hashicorp-vault-0.13.0/hashicorp-vault-0.13.0.tgz",
		"https://github.com/openshift-helm-charts/charts/releases/download/hashicorp-vault-0.12.0/hashicorp-vault-0.12.0.tgz",
		"https://github.com/IBM/charts/blob/master/repo/ibm-helm/ibm-object-storage-plugin-2.1.2.tgz?raw=true"}

	for _, chart := range charts {

		var out bytes.Buffer
		cmd1 := exec.Command("curl", "-L", chart)
		cmd2 := exec.Command("sha256sum")

		var err error
		cmd2.Stdin, err = cmd1.StdoutPipe()
		assert.NoError(t, err, "error get pipe from curl -L command")
		cmd2.Stdout = &out
		err = cmd2.Start()
		assert.NoError(t, err, "error starting sha256sum command")
		err = cmd1.Run()
		assert.NoError(t, err, "error starting curl -L command")
		err = cmd2.Wait()
		assert.NoError(t, err, "error waiting for sha256sum command")
		commandResponse := strings.Split(out.String(), " ")
		assert.Equal(t, commandResponse[0], GetPackageDigest(chart), fmt.Sprintf("%s digests did not match as expected", chart))

	}

}
