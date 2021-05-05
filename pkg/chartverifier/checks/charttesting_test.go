package checks

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"helm.sh/helm/v3/pkg/cli"
)

func TestChartTesting(t *testing.T) {
	type testCase struct {
		config      map[string]interface{}
		description string
		uri         string
	}

	testCases := []testCase{
		{
			config:      map[string]interface{}{},
			description: "with chart-testing defaults",
			uri:         "chart-0.1.0-v3.valid.tgz",
		},
		{
			config: map[string]interface{}{
				"upgrade": true,
			},
			description: "override chart-testing upgrade",
			uri:         "chart-0.1.0-v3.valid.tgz",
		},
		{
			config: map[string]interface{}{
				"skipMissingValues": true,
			},
			description: "override chart-testing upgrade",
			uri:         "chart-0.1.0-v3.valid.tgz",
		},
		{
			config: map[string]interface{}{
				"namespace": "ct-test-namespace",
			},
			description: "override chart-testing namespace",
			uri:         "chart-0.1.0-v3.valid.tgz",
		},
		{
			config: map[string]interface{}{
				"releaseLabel": "chart-verifier-app.kubernetes.io/instance",
			},
			description: "override chart-testing releaseLabel",
			uri:         "chart-0.1.0-v3.valid.tgz",
		},
	}

	for _, tc := range testCases {
		config := viper.New()
		settings := cli.New()

		_ = config.MergeConfigMap(tc.config)

		t.Run(tc.description, func(t *testing.T) {
			t.Skip()
			r, err := ChartTesting(
				&CheckOptions{
					URI:             tc.uri,
					ViperConfig:     config,
					HelmEnvSettings: settings,
				},
			)
			require.NoError(t, err)
			require.NotNil(t, r)
			require.True(t, r.Ok)
		})
	}
}
