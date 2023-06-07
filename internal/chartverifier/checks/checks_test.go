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
	"fmt"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"helm.sh/helm/v3/pkg/cli"

	"github.com/redhat-certification/chart-verifier/internal/chartverifier/pyxis"
	"github.com/redhat-certification/chart-verifier/internal/tool"
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
		settings := cli.New()
		t.Run(tc.description, func(t *testing.T) {
			r, err := IsHelmV3(&CheckOptions{URI: tc.uri, ViperConfig: config, HelmEnvSettings: settings})
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

func TestNotContainValuesSchemaRemoteRef(t *testing.T) {
	type testCase struct {
		description string
		uri         string
	}

	positiveTestCases := []testCase{
		{description: "chart with schema without remote ref", uri: "chart-0.1.0-v3.valid.tgz"},
	}

	for _, tc := range positiveTestCases {
		t.Run(tc.description, func(t *testing.T) {
			config := viper.New()
			r, err := NotContainValuesSchemaRemoteRef(&CheckOptions{URI: tc.uri, ViperConfig: config, HelmEnvSettings: cli.New()})
			require.NoError(t, err)
			require.NotNil(t, r)
			require.True(t, r.Ok)
			require.Equal(t, ValuesSchemaHasNoRemoteRef, r.Reason)
		})
	}

	negativeTestCases := []testCase{
		{description: "chart with schema with remote ref", uri: "chart-0.1.0-v3.schema-with-remote-ref.tgz"},
	}

	for _, tc := range negativeTestCases {
		t.Run(tc.description, func(t *testing.T) {
			config := viper.New()
			r, err := NotContainValuesSchemaRemoteRef(&CheckOptions{URI: tc.uri, ViperConfig: config, HelmEnvSettings: cli.New()})
			require.NoError(t, err)
			require.NotNil(t, r)
			require.False(t, r.Ok)
			require.Equal(t, ValuesSchemaHasRemoteRef, r.Reason)
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
		{description: "Helm lint works for chart with lint INFO reason", uri: "chart-0.1.0-v2.lint-info.tgz"},
		{description: "Helm lint works for chart with lint WARNING reason", uri: "chart-0.1.0-v2.lint-warning.tgz"},
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
	checkImages(t, ImagesAreCertified, false)

	checkImages(t, ImagesAreCertified_V1_1, true)
}

func checkImages(t *testing.T, fn func(*CheckOptions) (Result, error), version1_1 bool) {
	type testCase struct {
		description string
		uri         string
		numErrors   int
		numPasses   int
		numSkips    int
	}

	var testCases []testCase

	if !version1_1 {
		testCases = []testCase{
			{description: "chart-0.1.0-v3.valid.tgz check images passes", uri: "chart-0.1.0-v3.valid.tgz", numErrors: 0, numPasses: 5},
			{description: "Helm check images fails", uri: "chart-0.1.0-v3.with-crd.tgz", numErrors: 2, numPasses: 0},
			{description: "Helm check images fails", uri: "chart-0.1.0-v3.with-csi.tgz", numErrors: 1, numPasses: 0},
		}
	} else {
		testCases = []testCase{
			{description: "chart-0.1.0-v3.valid.tgz check images passes", uri: "chart-0.1.0-v3.valid.tgz", numErrors: 0, numPasses: 5, numSkips: 0},
			{description: "chart-0.1.0-v3.valid-skipped-images.tgz check images passes", uri: "chart-0.1.0-v3.valid-skipped-images.tgz", numErrors: 0, numPasses: 3, numSkips: 2},
			{description: "chart-0.1.0-v3.failed-skipped-images.tgz check images passes", uri: "chart-0.1.0-v3.failed-skipped-images.tgz", numErrors: 1, numPasses: 0, numSkips: 4},
			{description: "chart-0.1.0-v3.skipped-images.tgz check images passes", uri: "chart-0.1.0-v3.skipped-images.tgz", numErrors: 0, numPasses: 0, numSkips: 5},
			{description: "Helm check images fails", uri: "chart-0.1.0-v3.with-crd.tgz", numErrors: 2, numPasses: 0, numSkips: 0},
			{description: "Helm check images fails", uri: "chart-0.1.0-v3.with-csi.tgz", numErrors: 1, numPasses: 0, numSkips: 0},
		}
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			config := viper.New()
			r, err := fn(&CheckOptions{URI: tc.uri, ViperConfig: config, HelmEnvSettings: cli.New()})
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
			if tc.numSkips > 0 {
				for i := 0; i < tc.numSkips; i++ {
					require.Contains(t, r.Reason, ImageCertifySkipped)
					r.Reason = strings.Replace(r.Reason, ImageCertifySkipped, "_replaced_", 1)
				}
				require.False(t, strings.Contains(r.Reason, ImageCertifySkipped))
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
		{"Single repo Default version 1", "repo", &pyxis.ImageReference{Registries: []string(nil), Repository: "repo", Tag: "latest", Sha: ""}},
		{"Single repo Default version 2", "repo:", &pyxis.ImageReference{Registries: []string(nil), Repository: "repo", Tag: "latest", Sha: ""}},
		{"Single repo with version", "repo:1.1.8", &pyxis.ImageReference{Registries: []string(nil), Repository: "repo", Tag: "1.1.8", Sha: ""}},
		{"Double repo with version", "repo/product:1.1.8", &pyxis.ImageReference{Registries: []string(nil), Repository: "repo/product", Tag: "1.1.8", Sha: ""}},
		{"Triple repo with version", "repo/subrepo/product:1.1.8", &pyxis.ImageReference{Registries: []string(nil), Repository: "repo/subrepo/product", Tag: "1.1.8", Sha: ""}},
		{"Registry, single repo with version", "registry.com/product:1.1.8", &pyxis.ImageReference{Registries: []string{"registry.com"}, Repository: "product", Tag: "1.1.8", Sha: ""}},
		{"Registry, double repo with version", "registry.com/repo/product:1.1.8", &pyxis.ImageReference{Registries: []string{"registry.com"}, Repository: "repo/product", Tag: "1.1.8", Sha: ""}},
		{"Registry with port, double repo with version", "registry.com:8080/repo/product:1.1.8", &pyxis.ImageReference{Registries: []string{"registry.com:8080"}, Repository: "repo/product", Tag: "1.1.8", Sha: ""}},
		{"Single repo Sha256", "repo@sha256:12345", &pyxis.ImageReference{Registries: []string(nil), Repository: "repo", Tag: "", Sha: "sha256:12345"}},
		{"Single repo Sha128", "repo@sha128:12345", &pyxis.ImageReference{Registries: []string(nil), Repository: "repo", Tag: "", Sha: "sha128:12345"}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			imageRef := parseImageReference(testCase.image)
			require.Equal(t, *testCase.expectedImageRef, imageRef)
		})
	}
}

func TestRequiredAnnotationsPresent(t *testing.T) {
	type testCase struct {
		description string
		uri         string
	}

	positiveTestCases := []testCase{
		{description: "chart with no missing required annotations", uri: "chart-0.1.0-v3.no-missing-annotations.tgz"},
	}

	for _, tc := range positiveTestCases {
		t.Run(tc.description, func(t *testing.T) {
			config := viper.New()
			r, err := RequiredAnnotationsPresent(&CheckOptions{URI: tc.uri, ViperConfig: config, HelmEnvSettings: cli.New()})
			require.NoError(t, err)
			require.NotNil(t, r)
			require.True(t, r.Ok)
			require.Equal(t, RequiredAnnotationsSuccess, r.Reason)
		})
	}

	negativeTestCases := []testCase{
		{description: "chart with missing required annotations", uri: "chart-0.1.0-v3.missing-annotations.tgz"},
	}

	for _, tc := range negativeTestCases {
		t.Run(tc.description, func(t *testing.T) {
			message := fmt.Sprintf("%s: %v", RequiredAnnotationsFailure, requiredAnnotations)
			config := viper.New()
			r, err := RequiredAnnotationsPresent(&CheckOptions{URI: tc.uri, ViperConfig: config, HelmEnvSettings: cli.New()})
			require.NoError(t, err)
			require.NotNil(t, r)
			require.False(t, r.Ok)
			require.Equal(t, message, r.Reason)
		})
	}
}

func TestSemVers(t *testing.T) {
	// Vault: kubeVersion: '>= 1.14.0-0'
	// IBM: kubeversion: '>=1.10.1-0'
	// Fortanix : kubeversion: '>= 1.16.0 < 1.22.0'
	// Wildfly: kubeversion: ""
	// Infispan: kubeversion: ""

	type testCase struct {
		kubeVersion string
		OCPRange    string
	}

	// nolint:unused // TODO(komish): Identify historical purpose of this type
	// before considering for removal.
	type TestError struct {
		expectedOCPRange string
		gotOCORange      string
	}

	testCases := []testCase{
		{kubeVersion: "~1.22-0", OCPRange: "4.9"},
		{kubeVersion: "1.22.*", OCPRange: "4.9"},
		{kubeVersion: "^1.22", OCPRange: ">=4.9"},
		{kubeVersion: ">=1.20-0", OCPRange: ">=4.7"},
		{kubeVersion: "1.21 - 1.22", OCPRange: "4.8 - 4.9"},
		{kubeVersion: ">1.20", OCPRange: ">=4.8"},
		{kubeVersion: "~1.21", OCPRange: "4.8"},
		{kubeVersion: ">= 1.14.0-0", OCPRange: ">=4.2"},
		{kubeVersion: "1.16 - 1.21", OCPRange: "4.3 - 4.8"},
		{kubeVersion: "*", OCPRange: ">=4.1"},
		{kubeVersion: ">=1.16.0 <1.22.0", OCPRange: "Error converting kubeVersion to an OCP range : improper constraint: >=1.16.0 <1.22.0"},
	}

	for _, test := range testCases {
		t.Run(fmt.Sprintf("Check kube version %s", test.kubeVersion), func(t *testing.T) {
			OCPRange, err := getOCPRange(test.kubeVersion)
			if err != nil {
				require.Equal(t, test.OCPRange, fmt.Sprintf("%v", err))
			} else {
				require.Equal(t, test.OCPRange, OCPRange)
			}
		})
	}
}

func TestSignatureIsValid(t *testing.T) {
	type testCase struct {
		description string
		uri         string
		keyFile     string
		reason      string
		ok          bool
		skipped     bool
	}

	testCases := []testCase{
		{
			description: "unsigned chart",
			uri:         "chart-0.1.0-v3.no-missing-annotations.tgz",
			keyFile:     "",
			reason:      fmt.Sprintf("%s : %s", ChartNotSigned, SignatureIsNotPresentSuccess),
			ok:          true, skipped: true,
		},
		{
			description: "unsigned chart with key provided",
			uri:         "chart-0.1.0-v3.no-missing-annotations.tgz",
			keyFile:     "../../../tests/charts/psql-service/0.1.11/psql-service-0.1.11.tgz.key",
			reason:      fmt.Sprintf("%s : %s", ChartNotSigned, SignatureIsNotPresentSuccess),
			ok:          true, skipped: true,
		},

		{
			description: "signed chart with valid key",
			uri:         "../../../tests/charts/psql-service/0.1.11/psql-service-0.1.11.tgz",
			keyFile:     "../../../tests/charts/psql-service/0.1.11/psql-service-0.1.11.tgz.key",
			reason:      fmt.Sprintf("%s : %s", ChartSigned, SignatureIsValidSuccess),
			ok:          true, skipped: false,
		},
		{
			description: "signed chart with no key",
			uri:         "https://github.com/redhat-certification/chart-verifier/blob/main/tests/charts/psql-service/0.1.11/psql-service-0.1.11.tgz?raw=true",
			keyFile:     "",
			reason:      fmt.Sprintf("%s : %s", ChartSigned, SignatureNoKey),
			ok:          true, skipped: true,
		},
		{
			description: "signed chart with bad key",
			uri:         "../../../tests/charts/psql-service/0.1.11/psql-service-0.1.11.tgz",
			keyFile:     "../../../tests/charts/psql-service/0.1.11/psql-service-0.1.11.tgz.badkey",
			reason:      fmt.Sprintf("%s : %s", ChartSigned, SignatureFailure),
			ok:          false, skipped: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			config := viper.New()
			base64Key := ""
			var encodeErr error
			if len(tc.keyFile) > 0 {
				base64Key, encodeErr = tool.GetEncodedKey(tc.keyFile)
				require.NoError(t, encodeErr)
			}
			r, err := SignatureIsValid(&CheckOptions{URI: tc.uri, ViperConfig: config, HelmEnvSettings: cli.New(), PublicKeys: []string{base64Key}})
			require.NoError(t, err)
			require.NotNil(t, r)
			require.Equal(t, r.Ok, tc.ok, fmt.Sprintf("%s : outcome mismatch", tc.description))
			require.Equal(t, r.Skipped, tc.skipped, fmt.Sprintf("%s : skipped mismatch", tc.description))
			require.Contains(t, r.Reason, tc.reason, fmt.Sprintf("%s : reason mismatch", tc.description))
		})
	}
}
