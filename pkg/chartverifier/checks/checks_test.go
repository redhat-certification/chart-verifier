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
	"github.com/redhat-certification/chart-verifier/pkg/chartverifier/pyxis"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"helm.sh/helm/v3/pkg/cli"
)

func TestIsHelmV3(t *testing.T) {
	type testCase struct {
		description string
		uri         string
	}

	positiveTestCases := []testCase{
		{description: "valid tarball", uri: "chart-0.1.0-v3.valid.tgz"},
	}

	for _, tc := range positiveTestCases {
		config := viper.New()
		settings := cli.New()
		t.Run(tc.description, func(t *testing.T) {
			r, err := IsHelmV3(&CheckOptions{URI: tc.uri, ViperConfig: config, HelmEnvSettings: settings})
			require.NoError(t, err)
			require.NotNil(t, r)
			require.True(t, r.Ok)
			require.Equal(t, Helm3Reason, r.Reason)
		})
	}

	negativeTestCases := []testCase{
		{description: "invalid tarball", uri: "chart-0.1.0-v2.invalid.tgz"},
	}

	for _, tc := range negativeTestCases {
		config := viper.New()
		t.Run(tc.description, func(t *testing.T) {
			r, err := IsHelmV3(&CheckOptions{URI: tc.uri, ViperConfig: config})
			require.NoError(t, err)
			require.NotNil(t, r)
			require.False(t, r.Ok)
			require.Equal(t, NotHelm3Reason, r.Reason)
		})
	}
}

func TestHasReadme(t *testing.T) {
	type testCase struct {
		description string
		uri         string
	}

	positiveTestCases := []testCase{
		{description: "chart with README", uri: "chart-0.1.0-v3.valid.tgz"},
	}

	for _, tc := range positiveTestCases {
		t.Run(tc.description, func(t *testing.T) {
			config := viper.New()
			r, err := HasReadme(&CheckOptions{URI: tc.uri, ViperConfig: config, HelmEnvSettings: cli.New()})
			require.NoError(t, err)
			require.NotNil(t, r)
			require.True(t, r.Ok)
			require.Equal(t, ReadmeExist, r.Reason)
		})
	}

	negativeTestCases := []testCase{
		{description: "chart with README", uri: "chart-0.1.0-v3.without-readme.tgz"},
	}

	for _, tc := range negativeTestCases {
		t.Run(tc.description, func(t *testing.T) {
			config := viper.New()
			r, err := HasReadme(&CheckOptions{URI: tc.uri, ViperConfig: config, HelmEnvSettings: cli.New()})
			require.NoError(t, err)
			require.NotNil(t, r)
			require.False(t, r.Ok)
			require.Equal(t, ReadmeDoesNotExist, r.Reason)
		})
	}
}

func TestContainsTest(t *testing.T) {
	type testCase struct {
		description string
		uri         string
	}

	positiveTestCases := []testCase{
		{description: "tarball contains at least one test", uri: "chart-0.1.0-v3.valid.tgz"},
	}

	for _, tc := range positiveTestCases {
		t.Run(tc.description, func(t *testing.T) {
			config := viper.New()
			r, err := ContainsTest(&CheckOptions{URI: tc.uri, ViperConfig: config, HelmEnvSettings: cli.New()})
			require.NoError(t, err)
			require.NotNil(t, r)
			require.True(t, r.Ok)
			require.Equal(t, ChartTestFilesExist, r.Reason)
		})
	}

	negativeTestCases := []testCase{
		{description: "tarball contains at least one test", uri: "chart-0.1.0-v3.valid.notest.tgz"},
	}

	for _, tc := range negativeTestCases {
		t.Run(tc.description, func(t *testing.T) {
			config := viper.New()
			r, err := ContainsTest(&CheckOptions{URI: tc.uri, ViperConfig: config, HelmEnvSettings: cli.New()})
			require.NoError(t, err)
			require.NotNil(t, r)
			require.False(t, r.Ok)
			require.Equal(t, ChartTestFilesDoesNotExist, r.Reason)
		})
	}
}

func TestHasValuesSchema(t *testing.T) {
	type testCase struct {
		description string
		uri         string
	}

	positiveTestCases := []testCase{
		{description: "chart with values", uri: "chart-0.1.0-v3.valid.tgz"},
	}

	for _, tc := range positiveTestCases {
		t.Run(tc.description, func(t *testing.T) {
			config := viper.New()
			r, err := ContainsValuesSchema(&CheckOptions{URI: tc.uri, ViperConfig: config, HelmEnvSettings: cli.New()})
			require.NoError(t, err)
			require.NotNil(t, r)
			require.True(t, r.Ok)
			require.Equal(t, ValuesSchemaFileExist, r.Reason)
		})
	}

	negativeTestCases := []testCase{
		{description: "chart without values", uri: "chart-0.1.0-v3.no-values-schema.tgz"},
	}

	for _, tc := range negativeTestCases {
		t.Run(tc.description, func(t *testing.T) {
			config := viper.New()
			r, err := ContainsValuesSchema(&CheckOptions{URI: tc.uri, ViperConfig: config, HelmEnvSettings: cli.New()})
			require.NoError(t, err)
			require.NotNil(t, r)
			require.False(t, r.Ok)
			require.Equal(t, ValuesSchemaFileDoesNotExist, r.Reason)
		})
	}
}

func TestHasValues(t *testing.T) {
	type testCase struct {
		description string
		uri         string
	}

	positiveTestCases := []testCase{
		{description: "chart with values", uri: "chart-0.1.0-v3.valid.tgz"},
	}

	for _, tc := range positiveTestCases {
		t.Run(tc.description, func(t *testing.T) {
			config := viper.New()
			r, err := ContainsValues(&CheckOptions{URI: tc.uri, ViperConfig: config, HelmEnvSettings: cli.New()})
			require.NoError(t, err)
			require.NotNil(t, r)
			require.True(t, r.Ok)
			require.Equal(t, ValuesFileExist, r.Reason)
		})
	}

	negativeTestCases := []testCase{
		{description: "chart without values", uri: "chart-0.1.0-v3.no-values.tgz"},
	}

	for _, tc := range negativeTestCases {
		t.Run(tc.description, func(t *testing.T) {
			config := viper.New()
			r, err := ContainsValues(&CheckOptions{URI: tc.uri, ViperConfig: config, HelmEnvSettings: cli.New()})
			require.NoError(t, err)
			require.NotNil(t, r)
			require.False(t, r.Ok)
			require.Equal(t, ValuesFileDoesNotExist, r.Reason)
		})
	}
}

func TestHasMinKubeVersion(t *testing.T) {
	type testCase struct {
		description string
		uri         string
	}

	positiveTestCases := []testCase{
		{description: "minimum Kubernetes version specified", uri: "chart-0.1.0-v3.valid.tgz"},
	}

	for _, tc := range positiveTestCases {
		t.Run(tc.description, func(t *testing.T) {
			config := viper.New()
			r, err := HasKubeVersion(&CheckOptions{URI: tc.uri, ViperConfig: config, HelmEnvSettings: cli.New()})
			require.NoError(t, err)
			require.NotNil(t, r)
			require.True(t, r.Ok)
			require.Equal(t, KuberVersionSpecified, r.Reason)
		})
	}

	negativeTestCases := []testCase{
		{description: "minimum Kubernetes version not specified", uri: "chart-0.1.0-v3.without-minkubeversion.tgz"},
	}

	for _, tc := range negativeTestCases {
		t.Run(tc.description, func(t *testing.T) {
			config := viper.New()
			r, err := HasKubeVersion(&CheckOptions{URI: tc.uri, ViperConfig: config, HelmEnvSettings: cli.New()})
			require.NoError(t, err)
			require.NotNil(t, r)
			require.False(t, r.Ok)
			require.Equal(t, KuberVersionNotSpecified, r.Reason)
		})
	}

}

func TestNotContainCRDs(t *testing.T) {
	type testCase struct {
		description string
		uri         string
	}

	positiveTestCases := []testCase{
		{description: "Not contain CRDs", uri: "chart-0.1.0-v3.valid.tgz"},
	}

	for _, tc := range positiveTestCases {
		t.Run(tc.description, func(t *testing.T) {
			config := viper.New()
			r, err := NotContainCRDs(&CheckOptions{URI: tc.uri, ViperConfig: config, HelmEnvSettings: cli.New()})
			require.NoError(t, err)
			require.NotNil(t, r)
			require.True(t, r.Ok)
			require.Equal(t, ChartDoesNotContainCRDs, r.Reason)
		})
	}

	negativeTestCases := []testCase{
		{description: "Contain CRDs", uri: "chart-0.1.0-v3.with-crd.tgz"},
	}

	for _, tc := range negativeTestCases {
		t.Run(tc.description, func(t *testing.T) {
			config := viper.New()
			r, err := NotContainCRDs(&CheckOptions{URI: tc.uri, ViperConfig: config, HelmEnvSettings: cli.New()})
			require.NoError(t, err)
			require.NotNil(t, r)
			require.False(t, r.Ok)
			require.Equal(t, ChartContainCRDs, r.Reason)
		})
	}
}

func TestNotContainCSIObjects(t *testing.T) {
	type testCase struct {
		description string
		uri         string
	}

	positiveTestCases := []testCase{
		{description: "Not contain CSI objects", uri: "chart-0.1.0-v3.valid.tgz"},
	}

	for _, tc := range positiveTestCases {
		t.Run(tc.description, func(t *testing.T) {
			config := viper.New()
			r, err := NotContainCSIObjects(&CheckOptions{URI: tc.uri, ViperConfig: config, HelmEnvSettings: cli.New()})
			require.NoError(t, err)
			require.NotNil(t, r)
			require.True(t, r.Ok)
			require.Equal(t, CSIObjectsDoesNotExist, r.Reason)
		})
	}

	negativeTestCases := []testCase{
		{description: "Contain CRDs", uri: "chart-0.1.0-v3.with-csi.tgz"},
	}

	for _, tc := range negativeTestCases {
		t.Run(tc.description, func(t *testing.T) {
			config := viper.New()
			r, err := NotContainCSIObjects(&CheckOptions{URI: tc.uri, ViperConfig: config, HelmEnvSettings: cli.New()})
			require.NoError(t, err)
			require.NotNil(t, r)
			require.False(t, r.Ok)
			require.Equal(t, CSIObjectsExist, r.Reason)
		})
	}
}

func TestHelmLint(t *testing.T) {
	type testCase struct {
		description string
		uri         string
	}

	positiveTestCases := []testCase{
		{description: "Helm lint works for valid chart", uri: "chart-0.1.0-v3.valid.tgz"},
		{description: "Helm lint works for chart with lint INFO message", uri: "chart-0.1.0-v2.lint-info.tgz"},
		{description: "Helm lint works for chart with lint WARNING message", uri: "chart-0.1.0-v2.lint-warning.tgz"},
	}

	for _, tc := range positiveTestCases {
		t.Run(tc.description, func(t *testing.T) {
			config := viper.New()
			r, err := HelmLint(&CheckOptions{URI: tc.uri, ViperConfig: config, HelmEnvSettings: cli.New()})
			require.NoError(t, err)
			require.NotNil(t, r)
			require.True(t, r.Ok)
			require.Equal(t, HelmLintSuccessful, r.Reason)
		})
	}

	negativeTestCases := []testCase{
		{description: "Helm lint fails for chart with lint error", uri: "chart-0.1.0-v2.lint-error.tgz"},
	}

	for _, tc := range negativeTestCases {
		t.Run(tc.description, func(t *testing.T) {
			config := viper.New()
			r, err := HelmLint(&CheckOptions{URI: tc.uri, ViperConfig: config, HelmEnvSettings: cli.New()})
			require.NoError(t, err)
			require.NotNil(t, r)
			require.False(t, r.Ok)
			require.Contains(t, r.Reason, HelmLintHasFailedPrefix)
		})
	}

}

func TestImageCertify(t *testing.T) {

	type testCase struct {
		description string
		uri         string
		numErrors   int
		numPasses   int
	}

	testCases := []testCase{
		{description: "chart-0.1.0-v3.valid.tgz check images passes", uri: "chart-0.1.0-v3.valid.tgz", numErrors: 0, numPasses: 5},
		{description: "Helm check images fails", uri: "chart-0.1.0-v3.with-crd.tgz", numErrors: 2, numPasses: 0},
		{description: "Helm check images fails", uri: "chart-0.1.0-v3.with-csi.tgz", numErrors: 1, numPasses: 0},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			config := viper.New()
			r, err := ImagesAreCertified(&CheckOptions{URI: tc.uri, ViperConfig: config, HelmEnvSettings: cli.New()})
			require.NoError(t, err)
			require.NotNil(t, r)
			if tc.numErrors == 0 {
				require.True(t, r.Ok)
			} else {
				require.False(t, r.Ok)
				for i := 0; i < tc.numErrors; i++ {
					require.Contains(t, r.Reason, ImageNotCertified)
					r.Reason = strings.Replace(r.Reason, ImageNotCertified, "_replaced_", 1)
				}
				require.False(t, strings.Contains(r.Reason, ImageNotCertified))
			}
			if tc.numPasses > 0 {
				for i := 0; i < tc.numPasses; i++ {
					require.Contains(t, r.Reason, ImageCertified)
					r.Reason = strings.Replace(r.Reason, ImageCertified, "_replaced_", 1)
				}
				require.False(t, strings.Contains(r.Reason, ImageCertified))
			}
		})
	}

}

func TestImageParsing(t *testing.T) {

	type testCase struct {
		description      string
		image            string
		expectedImageRef *pyxis.ImageReference
	}

	testCases := []testCase{
		{"Single repo Default version 1", "repo", &pyxis.ImageReference{[]string(nil), "repo", "latest", ""}},
		{"Single repo Default version 2", "repo:", &pyxis.ImageReference{[]string(nil), "repo", "latest", ""}},
		{"Single repo with version", "repo:1.1.8", &pyxis.ImageReference{[]string(nil), "repo", "1.1.8", ""}},
		{"Double repo with version", "repo/product:1.1.8", &pyxis.ImageReference{[]string(nil), "repo/product", "1.1.8", ""}},
		{"Triple repo with version", "repo/subrepo/product:1.1.8", &pyxis.ImageReference{[]string(nil), "repo/subrepo/product", "1.1.8", ""}},
		{"Registry, single repo with version", "registry.com/product:1.1.8", &pyxis.ImageReference{[]string{"registry.com"}, "product", "1.1.8", ""}},
		{"Registry, double repo with version", "registry.com/repo/product:1.1.8", &pyxis.ImageReference{[]string{"registry.com"}, "repo/product", "1.1.8", ""}},
		{"Registry with port, double repo with version", "registry.com:8080/repo/product:1.1.8", &pyxis.ImageReference{[]string{"registry.com:8080"}, "repo/product", "1.1.8", ""}},
		{"Single repo Sha256", "repo@sha256:12345", &pyxis.ImageReference{[]string(nil), "repo", "", "sha256:12345"}},
		{"Single repo Sha128", "repo@sha128:12345", &pyxis.ImageReference{[]string(nil), "repo", "", "sha128:12345"}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			imageRef := parseImageReference(testCase.image)
			require.Equal(t, *testCase.expectedImageRef, imageRef)
		})
	}
}
