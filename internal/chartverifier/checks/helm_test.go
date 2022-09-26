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
	"io/ioutil"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"helm.sh/helm/v3/pkg/cli"

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
		{description: "chart-0.1.0-v3.valid.tgz images ", uri: "chart-0.1.0-v3.valid.tgz", images: []string{"registry.access.redhat.com/rhscl/postgresql-10-rhel7:1-66",
			"snyk/kubernetes-operator", "rhscl/mongodb-36-rhel7:1-65",
			"icr.io/cpopen/ibmcloud-object-storage-driver@sha256:fc17bb3e89d00b3eb0f50b3ea83aa75c52e43d8e56cf2e0f17475e934eeeeb5f",
			"icr.io/cpopen/ibmcloud-object-storage-plugin@sha256:cf654987c38d048bc9e654f3928e9ce9a2a4fd47ce0283bb5f339c1b99298e6e"}},
		{description: "chart-0.1.0-v3.with-crd.tgz", uri: "chart-0.1.0-v3.with-crd.tgz", images: []string{"nginx:1.16.0", "busybox"}},
		{description: "chart-0.1.0-v3.with-csi.tgz", uri: "chart-0.1.0-v3.with-csi.tgz", images: []string{"nginx:1.16.0"}},
	}

	for _, tc := range TestCases {
		t.Run(tc.description, func(t *testing.T) {
			images, err := getImageReferences(tc.uri, map[string]interface{}{})
			require.NoError(t, err)
			require.Equal(t, len(images), len(tc.images))
			for i := 0; i < len(tc.images); i++ {
				require.Contains(t, images, tc.images[i])
			}
		})
	}
}

func TestLongLineTemplate(t *testing.T) {

	content, err := ioutil.ReadFile("templates/test-template.yaml")
	require.NoError(t, err)

	images, err := getImagesFromContent(string(content))
	require.NoError(t, err)

	require.Equal(t, len(images), 2)

	require.Contains(t, images, "1.1.1/cv-test/image1:tag-123")
	require.Contains(t, images, "1.1.2/cv-test/image2:tag-223")

}
