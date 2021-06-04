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
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/redhat-certification/chart-verifier/pkg/testutil"
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

	ctx, cancel := context.WithCancel(context.Background())
	require.NoError(t, testutil.ServeCharts(ctx, addr, "./"))

	for _, tc := range positiveCases {
		t.Run(tc.description, func(t *testing.T) {
			c, _, err := LoadChartFromURI(tc.uri)
			require.NoError(t, err)
			require.NotNil(t, c)
		})
	}

	for _, tc := range negativeCases {
		t.Run(tc.description, func(t *testing.T) {
			c, _, err := LoadChartFromURI(tc.uri)
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
		{description: "chart-0.1.0-v3.valid.tgz images ", uri: "chart-0.1.0-v3.valid.tgz", images: []string{"nginx:latest",
			"snyk/kubernetes-operator", "rhscl/mongodb-36-rhel7:latest",
			"docker.io/ibmcom/ibmcloud-object-storage-driver@sha256:b6ec40ca7300bf9e2d0e7b9ff4272258f50d2d6ff9db766207f4a4281b2e33a1",
			"docker.io/ibmcom/ibmcloud-object-storage-plugin@sha256:0c361f70133a5aae4ac3cbbc250322f8dee2e71da734b818621033179508ce6f"}},
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
