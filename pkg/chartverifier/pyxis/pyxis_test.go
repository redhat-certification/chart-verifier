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

package pyxis

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_getImageRegistries(t *testing.T) {

	type testCase struct {
		description string
		repository  string
		registry    string
		message     string
	}

	PassTestCases := []testCase{
		{description: "Test nginx respository", repository: "nginx", registry: "registry.hub.docker.com", message: ""},
		{description: "Test rhel6.7 respository", repository: "rhel6.7", registry: "registry.access.redhat.com", message: ""},
		{description: "Test rhel8/nginx-116 respository", repository: "rhel8/nginx-116", registry: "registry.access.redhat.com", message: ""},
		{description: "Test ibm/nginx respository", repository: "ibm/nginx", registry: "non_registry", message: ""},
		{description: "Test turbonomic/zookeeper respository", repository: "turbonomic/zookeeper", registry: "registry.connect.redhat.com", message: ""}}

	for _, tc := range PassTestCases {
		t.Run(tc.description, func(t *testing.T) {
			reg, err := GetImageRegistries(tc.repository)
			require.NoError(t, err)
			require.Equal(t, tc.registry, reg[0])
		})
	}

	FailTestCases := []testCase{
		{description: "Test repository not found", repository: "ndoesnotexist", registry: "registry.hub.docker.com", message: "Respository not found"},
	}

	for _, tc := range FailTestCases {
		t.Run(tc.description, func(t *testing.T) {
			reg, err := GetImageRegistries(tc.repository)
			require.Error(t, err)
			require.Empty(t, reg)
			require.Contains(t, err.Error(), tc.message)
		})
	}

}

func Test_checkImageInRegistry(t *testing.T) {

	type testCase struct {
		description string
		repository  string
		registry    string
		version     string
		message     string
	}

	PassTestCases := []testCase{
		{description: "Test nginx registry and version found.", repository: "nginx", registry: "registry.hub.docker.com", version: "latest", message: ""},
		{description: "Test nginx rhel6.7 and version found.", repository: "rhel6.7", registry: "registry.access.redhat.com", version: "6.7", message: ""},
		{description: "Test rhel8/nginx-116 respository found.", repository: "rhel8/nginx-116", registry: "registry.access.redhat.com", version: "1-75", message: ""},
		{description: "Test turbonomic/zookeeper respository and version found.", repository: "turbonomic/zookeeper", registry: "registry.connect.redhat.com", version: "8.1.2", message: ""},
	}

	for _, tc := range PassTestCases {
		t.Run(tc.description, func(t *testing.T) {
			found, err := IsImageInRegistry(tc.repository, tc.version, tc.registry)
			require.NoError(t, err)
			require.True(t, found)
		})
	}

	FailTestCases := []testCase{
		{description: "Test nginx version not found", repository: "nginx", registry: "registry.hub.docker.com", version: "1.6.8", message: "Version 1.6.8 not found"},
		{description: "Test rhel6.7 registry not found", repository: "rhel6.7", registry: "registry.notfound.com", version: "7.8", message: "Registry not found: registry.notfound.com"},
	}

	for _, tc := range FailTestCases {
		t.Run(tc.description, func(t *testing.T) {
			found, err := IsImageInRegistry(tc.repository, tc.version, tc.registry)
			require.Error(t, err)
			require.False(t, found)
			require.Contains(t, err.Error(), tc.message)
		})
	}
}
