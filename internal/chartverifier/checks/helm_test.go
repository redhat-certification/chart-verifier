/*
 * Copyright 2021 Red Hat
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package checks

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"helm.sh/helm/v4/pkg/cli"

	"github.com/redhat-certification/chart-verifier/internal/testutil"
)

func TestLoadChartFromURI(t *testing.T) {
	addr := "127.0.0.1:9876"

	type testCase struct {
		description string
		uri         string
	}

	positiveCases := []testCase{
		{
			uri:         "chart-0.1.0-v3.valid.tgz",
			description: "absolute path",
		},
		{
			uri:         "http://" + addr + "/charts/chart-0.1.0-v3.valid.tgz",
			description: "remote path, http",
		},
	}

	negativeCases := []testCase{
		{
			uri:         "chart-0.1.0-v3.non-existing.tgz",
			description: "non existing file",
		},
		{
			uri:         "http://" + addr + "/charts/chart-0.1.0-v3.non-existing.tgz",
			description: "non existing remote file",
		},
	}

	repositoryCacheSetCases := []testCase{
		{
			uri:         "chart-0.1.0-v3.valid.tgz",
			description: "temporary cache defined",
		},
	}

	negativeCasesRepositoryCache := []testCase{
		{
			uri:         "chart-0.1.0-v3.non-existing.tgz",
			description: "non existing file with repository cache set",
		},
		{
			uri:         "http://" + addr + "/charts/chart-0.1.0-v3.non-existing.tgz",
			description: "non existing remote file with repository cache set",
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	require.NoError(t, testutil.ServeCharts(ctx, addr, "./"))

	for _, tc := range positiveCases {
		t.Run(tc.description, func(t *testing.T) {
			opts := CheckOptions{
				URI: tc.uri,
				Values: map[string]interface{}{
					"k8Project": "bogus",
				},
				ViperConfig:     viper.New(),
				HelmEnvSettings: cli.New(),
			}
			c, _, err := LoadChartFromURI(&opts)
			require.NoError(t, err)
			require.NotNil(t, c)
		})
	}

	for _, tc := range negativeCases {
		t.Run(tc.description, func(t *testing.T) {
			opts := CheckOptions{
				URI: tc.uri,
				Values: map[string]interface{}{
					"k8Project": "bogus",
				},
				ViperConfig:     viper.New(),
				HelmEnvSettings: cli.New(),
			}
			c, _, err := LoadChartFromURI(&opts)
			require.Error(t, err)
			require.True(t, IsChartNotFound(err))
			require.Equal(t, "chart not found: "+tc.uri, err.Error())
			require.Nil(t, c)
		})
	}

	for _, tc := range repositoryCacheSetCases {
		t.Run(tc.description, func(t *testing.T) {
			settings := cli.New()
			settings.RepositoryCache = "/tmp"
			opts := CheckOptions{
				URI:             tc.uri,
				ViperConfig:     viper.New(),
				HelmEnvSettings: settings,
			}
			c, _, err := LoadChartFromURI(&opts)
			require.NoError(t, err)
			require.NotNil(t, c)
		})
	}
	for _, tc := range negativeCasesRepositoryCache {
		t.Run(tc.description, func(t *testing.T) {
			settings := cli.New()
			settings.RepositoryCache = "/tmp"
			opts := CheckOptions{
				URI:             tc.uri,
				ViperConfig:     viper.New(),
				HelmEnvSettings: settings,
			}
			c, _, err := LoadChartFromURI(&opts)
			require.Error(t, err)
			require.True(t, IsChartNotFound(err))
			require.Equal(t, "chart not found: "+tc.uri, err.Error())
			require.Nil(t, c)
		})
	}
	cancel()
}

func TestTemplate(t *testing.T) {
	type testCase struct {
		description string
		uri         string
		images      []string
	}

	TestCases := []testCase{
		{description: "chart-0.1.0-v3.valid.tgz images ", uri: "chart-0.1.0-v3.valid.tgz", images: []string{
			"registry.access.redhat.com/rhscl/postgresql-10-rhel7:1-161",
		}},
		{description: "chart-0.1.0-v3.with-crd.tgz", uri: "chart-0.1.0-v3.with-crd.tgz", images: []string{"nginx:1.16.0"}},
		{description: "chart-0.1.0-v3.with-csi.tgz", uri: "chart-0.1.0-v3.with-csi.tgz", images: []string{"nginx:1.16.0"}},
	}

	for _, tc := range TestCases {
		t.Run(tc.description, func(t *testing.T) {
			images, err := getImageReferences(tc.uri, map[string]interface{}{}, defaultMockedKubeVersionString)
			require.NoError(t, err)
			require.Equal(t, len(images), len(tc.images))
			for i := 0; i < len(tc.images); i++ {
				require.Contains(t, images, tc.images[i])
			}
		})
	}
}

func TestLongLineTemplate(t *testing.T) {
	content, err := os.ReadFile("templates/test-template.yaml")
	require.NoError(t, err)

	images, err := getImagesFromContent(string(content))
	require.NoError(t, err)

	require.Equal(t, len(images), 2)

	require.Contains(t, images, "1.1.1/cv-test/image1:tag-123")
	require.Contains(t, images, "1.1.2/cv-test/image2:tag-223")
}

func TestGetImagesFromContent(t *testing.T) {
	testCases := []struct {
		name    string
		content string
		want    []string
	}{
		{
			name: "find images in yaml",
			content: `
	image: "registry.access.redhat.com/rhscl/postgresql-10-rhel7:1-161"
	image: 'busybox'
	image:  "  "
	image: registry.redhat.io/cpopen/ibmcloud-object-storage-driver@sha256:fc17bb3e89d00b3eb0f50b3ea83aa75c52e43d8e56cf2e0f17475e934eeeeb5f
`,
			want: []string{
				"registry.access.redhat.com/rhscl/postgresql-10-rhel7:1-161",
				"busybox",
				"",
				"registry.redhat.io/cpopen/ibmcloud-object-storage-driver@sha256:fc17bb3e89d00b3eb0f50b3ea83aa75c52e43d8e56cf2e0f17475e934eeeeb5f",
			},
		},
		{
			name: "do not match against mappings",
			content: `
			image:
			  repository: "registry.access.redhat.com/rhscl/postgresql-10-rhel7:1-161"
			`,
			want: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := getImagesFromContent(tc.content)
			require.Nil(t, err)
			if testing.Verbose() {
				t.Logf("got %d images", len(got))
				t.Logf("got: %s", got)
			}
			if len(got) != len(tc.want) {
				t.Errorf("got %d images but, want %d", len(got), len(tc.want))
			}
			for _, image := range got {
				if strings.TrimSpace(image) == "" {
					t.Logf("Found empty image")
				}
			}
		})
	}
}

func TestExcludeTestTemplates(t *testing.T) {
	testCases := []struct {
		name           string
		content        string
		expectedImages int
	}{
		{
			name: "excludes images from test templates identified by source path",
			content: `---
# Source: mychart/templates/deployment.yaml
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      containers:
        - image: registry.redhat.io/rhel8/nginx-120:latest
---
# Source: mychart/templates/tests/test-connection.yaml
apiVersion: v1
kind: Pod
spec:
  containers:
    - image: dotnet-runtime:latest
`,
			expectedImages: 1,
		},
		{
			name: "excludes images from test templates identified by hook annotation (double quotes)",
			content: `---
# Source: mychart/templates/deployment.yaml
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      containers:
        - image: registry.redhat.io/rhel8/nginx-120:latest
---
# Source: mychart/templates/test-pod.yaml
apiVersion: v1
kind: Pod
metadata:
  annotations:
    "helm.sh/hook": test
spec:
  containers:
    - image: busybox:latest
`,
			expectedImages: 1,
		},
		{
			name: "excludes images from test templates identified by hook annotation (no quotes)",
			content: `---
# Source: mychart/templates/hook-test.yaml
apiVersion: v1
kind: Pod
metadata:
  annotations:
    helm.sh/hook: test
spec:
  containers:
    - image: uncertified:latest
---
# Source: mychart/templates/deployment.yaml
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      containers:
        - image: registry.redhat.io/app:1.0
`,
			expectedImages: 1,
		},
		{
			name: "keeps all images when no test templates present",
			content: `---
# Source: mychart/templates/deployment.yaml
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      containers:
        - image: registry.redhat.io/app:1.0
---
# Source: mychart/templates/sidecar.yaml
apiVersion: v1
kind: Pod
spec:
  containers:
    - image: registry.redhat.io/sidecar:2.0
`,
			expectedImages: 2,
		},
		{
			name: "handles multiple test templates",
			content: `---
# Source: mychart/templates/deployment.yaml
spec:
  containers:
    - image: registry.redhat.io/app:1.0
---
# Source: mychart/templates/tests/test-1.yaml
spec:
  containers:
    - image: test-image-1:latest
---
# Source: mychart/templates/tests/test-2.yaml
spec:
  containers:
    - image: test-image-2:latest
`,
			expectedImages: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			filtered := excludeTestTemplates(tc.content)
			images, err := getImagesFromContent(filtered)
			require.Nil(t, err)
			if len(images) != tc.expectedImages {
				t.Errorf("got %d images, want %d. Images: %v", len(images), tc.expectedImages, images)
			}
		})
	}
}

func TestIsTestTemplate(t *testing.T) {
	testCases := []struct {
		name     string
		doc      string
		expected bool
	}{
		{
			name:     "source path with /tests/ directory",
			doc:      "# Source: mychart/templates/tests/test-connection.yaml\napiVersion: v1\nkind: Pod\n",
			expected: true,
		},
		{
			name:     "helm.sh/hook test annotation double-quoted",
			doc:      "apiVersion: v1\nmetadata:\n  annotations:\n    \"helm.sh/hook\": test\n",
			expected: true,
		},
		{
			name:     "helm.sh/hook test annotation single-quoted",
			doc:      "apiVersion: v1\nmetadata:\n  annotations:\n    'helm.sh/hook': test\n",
			expected: true,
		},
		{
			name:     "helm.sh/hook test annotation unquoted",
			doc:      "apiVersion: v1\nmetadata:\n  annotations:\n    helm.sh/hook: test\n",
			expected: true,
		},
		{
			name:     "regular deployment template",
			doc:      "# Source: mychart/templates/deployment.yaml\napiVersion: apps/v1\nkind: Deployment\n",
			expected: false,
		},
		{
			name:     "empty document",
			doc:      "",
			expected: false,
		},
		{
			name:     "path containing test but not in /tests/ directory",
			doc:      "# Source: mychart/templates/deployment-test-config.yaml\napiVersion: v1\n",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := isTestTemplate(tc.doc)
			if got != tc.expected {
				t.Errorf("isTestTemplate() = %v, want %v", got, tc.expected)
			}
		})
	}
}
