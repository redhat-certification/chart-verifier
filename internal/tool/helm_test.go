package tool

import (
	"context"
	"io"
	"os"
	"testing"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/stretchr/testify/require"
	"helm.sh/helm/v4/pkg/action"
	chartcommon "helm.sh/helm/v4/pkg/chart/common"
	chartv2 "helm.sh/helm/v4/pkg/chart/v2"
	"helm.sh/helm/v4/pkg/cli"
	kubefake "helm.sh/helm/v4/pkg/kube/fake"

	releasecommon "helm.sh/helm/v4/pkg/release/common"
	releasev1 "helm.sh/helm/v4/pkg/release/v1"
	"helm.sh/helm/v4/pkg/storage"
	"helm.sh/helm/v4/pkg/storage/driver"
)

func TestInstall(t *testing.T) {
	tests := []struct {
		releaseName string
		chartPath   string
		expected    string
		timeout     time.Duration
	}{
		{
			releaseName: "valid-chart",
			chartPath:   "../chartverifier/checks/psql-service-0.1.7",
			expected:    "",
			timeout:     10 * time.Second,
		},
		{
			releaseName: "invalid-chart",
			chartPath:   "../chartverifier/checks/psql-service-9.9.9",
			expected:    "path \"../chartverifier/checks/psql-service-9.9.9\" not found",
			timeout:     10 * time.Second,
		},
	}
	for _, tt := range tests {
		t.Run(tt.releaseName, func(t *testing.T) {
			actionConfig := &action.Configuration{
				Releases:     storage.Init(driver.NewMemory()),
				KubeClient:   &kubefake.PrintingKubeClient{Out: io.Discard},
				Capabilities: chartcommon.DefaultCapabilities,
			}
			helm := Helm{
				config:      actionConfig,
				args:        map[string]interface{}{"set": "k8Project=default"},
				envSettings: &cli.EnvSettings{},
			}

			ctx, cancel := context.WithTimeout(context.Background(), tt.timeout)
			defer cancel()

			beforeInstallTime := time.Now()
			err := helm.Install(ctx, "default", tt.chartPath, tt.releaseName, "")
			require.WithinDuration(t, beforeInstallTime, time.Now(), tt.timeout)
			if err == nil {
				require.Equal(t, tt.expected, "")
			} else {
				require.Equal(t, tt.expected, err.Error())
			}
		})
	}
}

func TestUninstall(t *testing.T) {
	tests := []struct {
		name     string
		release  *releasev1.Release
		expected string
	}{
		{
			name: "successful release uninstall should remove release installed",
			release: &releasev1.Release{
				Name: "test-release-valid",
				Info: &releasev1.Info{
					Status: releasecommon.StatusDeployed,
				},
				Namespace: "default",
			},
			expected: "",
		},
		{
			name: "remove non-existent release should result in error",
			release: &releasev1.Release{
				Name: "test-release-invalid",
				Info: &releasev1.Info{
					Status: releasecommon.StatusDeployed,
				},
				Namespace: "default",
			},
			expected: "uninstall: Release not loaded: test-release-invalid: release: not found",
		},
	}

	for _, tt := range tests {
		store := storage.Init(driver.NewMemory())
		t.Run(tt.name, func(t *testing.T) {
			actionConfig := &action.Configuration{
				Releases:     store,
				KubeClient:   &kubefake.PrintingKubeClient{Out: io.Discard},
				Capabilities: chartcommon.DefaultCapabilities,
			}
			helm := Helm{
				config:      actionConfig,
				args:        make(map[string]interface{}),
				envSettings: &cli.EnvSettings{},
			}
			if tt.expected == "" {
				// create fake release
				if err := store.Create(tt.release); err != nil {
					t.Error(err)
				}
			}
			err := helm.Uninstall("default", tt.release.Name)
			if err == nil {
				require.Equal(t, tt.expected, "")
			} else {
				require.Equal(t, tt.expected, err.Error())
			}
		})
	}
}

func TestUpgrade(t *testing.T) {
	testValues, err := chartcommon.ReadValuesFile("../chartverifier/checks/psql-service-0.1.7/values.yaml")
	if err != nil {
		t.Error(err)
	}
	testValues.AsMap()["k8Project"] = "default"

	var chartMetadata chartv2.Metadata
	yamlFile, err := os.ReadFile("../chartverifier/checks/psql-service-0.1.7/Chart.yaml")
	require.NoError(t, err)
	err = yaml.Unmarshal(yamlFile, &chartMetadata)
	require.NoError(t, err)

	tests := []struct {
		name      string
		chartPath string
		release   *releasev1.Release
		expected  string
		timeout   time.Duration
	}{
		{
			name:      "successful release upgrade should not return error",
			chartPath: "../chartverifier/checks/psql-service-0.1.7",
			release: &releasev1.Release{
				Name: "test-release-valid",
				Info: &releasev1.Info{
					Status: releasecommon.StatusDeployed,
				},
				Namespace: "default",
				Chart:     &chartv2.Chart{Metadata: &chartMetadata, Values: testValues},
			},
			expected: "",
			timeout:  10 * time.Second,
		},
		{
			name:      "upgrade non-existent release should result in error",
			chartPath: "../chartverifier/checks/psql-service-0.1.7",
			release: &releasev1.Release{
				Name: "test-release-invalid",
				Info: &releasev1.Info{
					Status: releasecommon.StatusDeployed,
				},
				Namespace: "default",
				Chart:     &chartv2.Chart{Values: testValues},
			},
			expected: "\"test-release-invalid\" has no deployed releases",
			timeout:  10 * time.Second,
		},
	}

	for _, tt := range tests {
		store := storage.Init(driver.NewMemory())
		t.Run(tt.name, func(t *testing.T) {
			actionConfig := &action.Configuration{
				Releases:     store,
				KubeClient:   &kubefake.PrintingKubeClient{Out: io.Discard},
				Capabilities: chartcommon.DefaultCapabilities,
			}
			helm := Helm{
				config:      actionConfig,
				args:        make(map[string]interface{}),
				envSettings: &cli.EnvSettings{},
			}
			if tt.expected == "" {
				// create fake release
				if err := store.Create(tt.release); err != nil {
					t.Error(err)
				}
			}
			ctx, cancel := context.WithTimeout(context.Background(), tt.timeout)
			defer cancel()
			beforeUpgradeTime := time.Now()
			err := helm.Upgrade(ctx, "default", tt.chartPath, tt.release.Name)
			require.WithinDuration(t, beforeUpgradeTime, time.Now(), tt.timeout)
			if err == nil {
				require.Equal(t, tt.expected, "")
			} else {
				require.Equal(t, tt.expected, err.Error())
			}
		})
	}
}

func TestReleaseTesting(t *testing.T) {
	releaseTestPath := "../chartverifier/checks/psql-service-0.1.7/templates/tests/test-psql-connection.yaml"
	releaseTest, err := os.ReadFile(releaseTestPath)
	if err != nil {
		t.Error(err)
	}
	testHooks := []*releasev1.Hook{
		{
			Name:     "test-success-hook",
			Kind:     "Pod",
			Path:     releaseTestPath,
			Manifest: string(releaseTest),
			LastRun:  releasev1.HookExecution{},
			Events:   []releasev1.HookEvent{releasev1.HookTest},
		},
	}
	tests := []struct {
		name      string
		chartPath string
		release   *releasev1.Release
		expected  string
		timeout   time.Duration
	}{
		{
			name:      "successful release test should not return error",
			chartPath: "../chartverifier/checks/psql-service-0.1.7",
			release: &releasev1.Release{
				Name: "test-release-valid",
				Info: &releasev1.Info{
					Status: releasecommon.StatusDeployed,
				},
				Namespace: "default",
				Hooks:     testHooks,
			},
			expected: "",
			timeout:  10 * time.Second,
		},
		{
			name:      "release test on non-existent release should result in error",
			chartPath: "../chartverifier/checks/psql-service-0.1.7",
			release: &releasev1.Release{
				Name: "test-release-invalid",
				Info: &releasev1.Info{
					Status: releasecommon.StatusDeployed,
				},
				Namespace: "default",
				Hooks:     testHooks,
			},
			expected: "release: not found",
			timeout:  10 * time.Second,
		},
		{
			name:      "release test with a zero or negative value should result in error",
			chartPath: "../chartverifier/checks/psql-service-0.1.7",
			release: &releasev1.Release{
				Name: "test-release-invalid-timeout",
				Info: &releasev1.Info{
					Status: releasecommon.StatusDeployed,
				},
				Namespace: "default",
				Hooks:     testHooks,
			},
			expected: "Helm test error : timeout has expired, please consider increasing the timeout using the chart-verifier timeout flag",
			timeout:  -1 * time.Second,
		},
	}

	for _, tt := range tests {
		store := storage.Init(driver.NewMemory())
		t.Run(tt.name, func(t *testing.T) {
			actionConfig := &action.Configuration{
				Releases:     store,
				KubeClient:   &kubefake.PrintingKubeClient{Out: io.Discard},
				Capabilities: chartcommon.DefaultCapabilities,
			}
			helm := Helm{
				config:      actionConfig,
				args:        make(map[string]interface{}),
				envSettings: &cli.EnvSettings{},
			}
			if tt.expected == "" {
				// create fake release
				if err := store.Create(tt.release); err != nil {
					t.Error(err)
				}
			}
			ctx, cancel := context.WithTimeout(context.Background(), tt.timeout)
			defer cancel()
			beforeTestime := time.Now()
			err := helm.Test(ctx, "default", tt.release.Name)
			if tt.timeout <= 0 {
				require.WithinDuration(t, beforeTestime, time.Now(), 1*time.Second)
			} else {
				require.WithinDuration(t, beforeTestime, time.Now(), tt.timeout)
			}
			if err == nil {
				require.Equal(t, tt.expected, "")
			} else {
				require.Equal(t, tt.expected, err.Error())
			}
		})
	}
}
