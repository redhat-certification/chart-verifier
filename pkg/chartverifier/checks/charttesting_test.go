package checks

import (
	"fmt"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"helm.sh/helm/v3/pkg/cli"
)

// absPathFromSourceFileLocation returns the absolute path of a file or directory under the current source file's
// location.
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

func lookPath(programs ...string) error {
	for _, p := range programs {
		_, err := exec.LookPath(p)
		if err != nil {
			return fmt.Errorf("required program %q not found", p)
		}
	}
	return nil
}

func TestChartTesting(t *testing.T) {
	if err := lookPath("helm", "kubectl"); err != nil {
		t.Skip(err.Error())
	}

	type testCase struct {
		description string
		opts        CheckOptions
	}

	chartUri, err := absPathFromSourceFileLocation("psql-service-0.1.7")
	if err != nil {
		t.Error(err)
	}

	positiveTestCases := []testCase{
		{
			description: "providing a valid k8Project value should succeed",
			opts: CheckOptions{
				URI: chartUri,
				Values: map[string]interface{}{
					"k8Project": "default",
				},
				ViperConfig:     viper.New(),
				HelmEnvSettings: cli.New(),
			},
		},
	}

	for _, tc := range positiveTestCases {
		t.Run(tc.description, func(t *testing.T) {
			r, err := ChartTesting(&tc.opts)
			require.NoError(t, err)
			require.NotNil(t, r)
			require.Equal(t, ChartTestingSuccess, r.Reason)
			require.True(t, r.Ok)
		})
	}

	negativeTestCases := []testCase{
		{
			description: "providing a bogus k8Project should fail",
			opts: CheckOptions{
				URI: chartUri,
				Values: map[string]interface{}{
					"k8Project": "bogus",
				},
				ViperConfig:     viper.New(),
				HelmEnvSettings: cli.New(),
			},
		},
	}

	for _, tc := range negativeTestCases {
		t.Run(tc.description, func(t *testing.T) {
			r, err := ChartTesting(&tc.opts)
			require.NotNil(t, r)
			require.False(t, r.Ok)
			require.NoError(t, err)
			require.Contains(t, r.Reason, "executing helm with args")
		})
	}
}
