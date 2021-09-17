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
		message     string
		imageRef    ImageReference
	}

	PassTestCases := []testCase{
		{description: "Test nginx registry and version found.", message: "", imageRef: ImageReference{Repository: "rhel6.7", Registries: []string{"registry.access.redhat.com"}, Tag: "6.7", Sha: ""}},
		{description: "Test nginx rhel6.7 and version found.", imageRef: ImageReference{Repository: "rhel6.7", Registries: []string{"registry.access.redhat.com"}, Tag: "6.7", Sha: ""}, message: ""},
		{description: "Test rhel8/nginx-116 respository found.", imageRef: ImageReference{Repository: "rhel8/nginx-116", Registries: []string{"registry.access.redhat.com"}, Tag: "1-75", Sha: ""}, message: ""},
		{description: "Test turbonomic/zookeeper respository and version found.", imageRef: ImageReference{Repository: "turbonomic/zookeeper", Registries: []string{"registry.connect.redhat.com"}, Tag: "8.1.2", Sha: ""}, message: ""},
		{description: "Test ibmcom/ibmcloud-object-storage-driver respository and sha found.", imageRef: ImageReference{Repository: "ibmcom/ibmcloud-object-storage-driver", Registries: []string{"docker.io"}, Tag: "", Sha: "sha256:e9152b9e7dfca10cf02f8de7a8d14a4067e7fe69695699411cadaf282263f099"}, message: ""},
		{description: "Test ibmcom/ibmcloud-object-storage-plugin respository and sha found.", imageRef: ImageReference{Repository: "ibmcom/ibmcloud-object-storage-plugin", Registries: []string{"docker.io"}, Tag: "", Sha: "sha256:b2e7a3ced38cb9197e7d1a3bea3ffcc9dda46e7d12ca5337e5ba8d4253659309"}, message: ""},
	}

	for _, tc := range PassTestCases {
		t.Run(tc.description, func(t *testing.T) {
			found, err := IsImageInRegistry(tc.imageRef)
			require.NoError(t, err)
			require.True(t, found)
		})
	}

	FailTestCases := []testCase{
		{description: "Test nginx version not found", imageRef: ImageReference{Repository: "nginx", Registries: []string{"registry.hub.docker.com"}, Tag: "1.6.8", Sha: ""}, message: "Tag 1.6.8 not found"},
		{description: "Test rhel6.7 registry not found", imageRef: ImageReference{Repository: "rhel6.7", Registries: []string{"registry.notfound.com"}, Tag: "7.8", Sha: ""}, message: "No images found for Registry/Repository: registry.notfound.com/rhel6.7"},
		{description: "Test ibmcom/ibmcloud-object-storage-plugin respository sha not found.", imageRef: ImageReference{Repository: "ibmcom/ibmcloud-object-storage-plugin", Registries: []string{"docker.io"}, Tag: "", Sha: "sha256:0d561f70133a5aae4ac3cbbc250322f8dee2e71da734b818621033179508ce6f"}, message: "Digest sha256:0d561f70133a5aae4ac3cbbc250322f8dee2e71da734b818621033179508ce6f not found"},
	}

	for _, tc := range FailTestCases {
		t.Run(tc.description, func(t *testing.T) {
			found, err := IsImageInRegistry(tc.imageRef)
			require.Error(t, err)
			require.False(t, found)
			require.Contains(t, err.Error(), tc.message)
		})
	}
}
