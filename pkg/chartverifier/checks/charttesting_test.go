package checks

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestChartTesting(t *testing.T) {
	type testCase struct {
		config      map[string]interface{}
		description string
		uri         string
	}

	testCases := []testCase{
		{
			config: map[string]interface{}{
				"buildId":           "[BUILD_ID]",
				"upgrade":           true,
				"skipMissingValues": true,
				"namespace":         "default",
				"releaseLabel":      "app.kubernetes.io/instance",
			},
			description: "valid tarball",
			uri:         "chart-0.1.0-v3.valid.tgz",
		},
	}

	for _, tc := range testCases {
		config := viper.New()
		_ = config.MergeConfigMap(tc.config)

		t.Run(tc.description, func(t *testing.T) {
			r, err := ChartTesting(&CheckOptions{URI: tc.uri, ViperConfig: config})
			require.NoError(t, err)
			require.NotNil(t, r)
			require.True(t, r.Ok)
		})
	}
}
