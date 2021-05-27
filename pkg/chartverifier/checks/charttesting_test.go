package checks

import (
	"fmt"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
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

func TestChartTesting(t *testing.T) {

	type testCase struct {
		description string
		opts        CheckOptions
	}

    chartUri, err := absPathFromSourceFileLocation("chart-0.1.0-v3.valid.tgz")
    if err != nil {
        t.Error(err)
    }

	positiveTestCases := []testCase{
		{
			description: "with license=true value override",
			opts: CheckOptions{
				URI: chartUri,
				Values: map[string]interface{}{
					"license": true,
				},
			},
		},
	}

	for _, tc := range positiveTestCases {
		t.Run(tc.description, func(t *testing.T) {
			r, err := ChartTesting(&tc.opts)
			require.NoError(t, err)
			require.NotNil(t, r)
			require.True(t, r.Ok)
		})
	}

	negativeTestCases := []testCase{
		{
			description: "with chart-testing defaults",
			opts: CheckOptions{
				URI: chartUri,
				Values: map[string]interface{}{},
			},
		},
	}

	for _, tc := range negativeTestCases {
		t.Run(tc.description, func(t *testing.T) {
			r, err := ChartTesting(&tc.opts)
			require.Error(t, err)
			require.NotNil(t, r)
			require.False(t, r.Ok)
		})
	}
}
