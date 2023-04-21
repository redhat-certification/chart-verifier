package checks

import (
	"errors"
	"fmt"
	"os"
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

func TestChartTesting(t *testing.T) {
	if os.Getenv("CHART_VERIFIER_ENABLE_CLUSTER_TESTING") == "" {
		t.Skip("CHART_VERIFIER_ENABLE_CLUSTER_TESTING not set, skipping in cluster tests")
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
		{
			// the chart being used in this test forces the rendered resources to have an empty namespace field, which
			// is invalid and can't be overriden using helm's namespace option.
			description: "empty values should fail",
			opts: CheckOptions{
				URI:             chartUri,
				Values:          map[string]interface{}{},
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

func getVersionError(settings *cli.EnvSettings) (string, error) {
	return "", errors.New("error")
}

func getVersionGood(settings *cli.EnvSettings) (string, error) {
	return "4.7.9", nil
}

type testAnnotationHolder struct {
	OpenShiftVersion              string
	CertifiedOpenShiftVersionFlag string
}

func (holder *testAnnotationHolder) SetCertifiedOpenShiftVersion(version string) {
	holder.OpenShiftVersion = version
}

func (holder *testAnnotationHolder) GetCertifiedOpenShiftVersionFlag() string {
	return holder.CertifiedOpenShiftVersionFlag
}

func (holder *testAnnotationHolder) SetSupportedOpenShiftVersions(version string) {}

func TestVersionSetting(t *testing.T) {
	type testCase struct {
		description string
		opts        *CheckOptions
		versioner   Versioner
		version     string
		error       string
	}

	testCases := []testCase{
		{
			description: "oc.Version returns 4.7.9",
			opts:        &CheckOptions{AnnotationHolder: &testAnnotationHolder{}},
			versioner:   getVersionGood,
			version:     "4.7.9",
		},
		{
			description: "oc.Version returns error, flag set to 4.7.8",
			opts:        &CheckOptions{AnnotationHolder: &testAnnotationHolder{CertifiedOpenShiftVersionFlag: "4.7.8"}},
			versioner:   getVersionError,
			version:     "4.7.8",
		},
		{
			description: "oc.Version returns semantic error, flag set to fourseveneight",
			opts:        &CheckOptions{AnnotationHolder: &testAnnotationHolder{CertifiedOpenShiftVersionFlag: "fourseveneight"}},
			versioner:   getVersionError,
			error:       "OpenShift version is not following SemVer spec. Invalid Semantic Version",
		},
		{
			description: "oc.Version returns error, flag not set",
			opts:        &CheckOptions{AnnotationHolder: &testAnnotationHolder{}},
			versioner:   getVersionError,
			error:       "Missing OpenShift version. error. And the 'openshift-version' flag has not set.",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			err := setOCVersion(tc.opts.AnnotationHolder, tc.opts.HelmEnvSettings, tc.versioner)

			if len(tc.error) > 0 {
				require.Error(t, err)
				require.Equal(t, tc.error, err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.version, tc.opts.AnnotationHolder.(*testAnnotationHolder).OpenShiftVersion)
			}
		})
	}
}
