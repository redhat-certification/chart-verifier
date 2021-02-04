/*
 * Copyright (C) 08/01/2021, 02:01, igors
 * This file is part of helmcertifier.
 *
 * helmcertifier is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * helmcertifier is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with helmcertifier.  If not, see <http://www.gnu.org/licenses/>.
 */

package checks

import (
	"testing"

	"github.com/stretchr/testify/require"
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
		t.Run(tc.description, func(t *testing.T) {
			r, err := IsHelmV3(tc.uri)
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
		t.Run(tc.description, func(t *testing.T) {
			r, err := IsHelmV3(tc.uri)
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
			r, err := HasReadme(tc.uri)
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
			r, err := HasReadme(tc.uri)
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
			r, err := ContainsTest(tc.uri)
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
			r, err := ContainsTest(tc.uri)
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
			r, err := ContainsValuesSchema(tc.uri)
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
			r, err := ContainsValuesSchema(tc.uri)
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
			r, err := ContainsValues(tc.uri)
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
			r, err := ContainsValues(tc.uri)
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
			r, err := HasMinKubeVersion(tc.uri)
			require.NoError(t, err)
			require.NotNil(t, r)
			require.True(t, r.Ok)
			require.Equal(t, MinKuberVersionSpecified, r.Reason)
		})
	}

	negativeTestCases := []testCase{
		{description: "minimum Kubernetes version not specified", uri: "chart-0.1.0-v3.without-minkubeversion.tgz"},
	}

	for _, tc := range negativeTestCases {
		t.Run(tc.description, func(t *testing.T) {
			r, err := HasMinKubeVersion(tc.uri)
			require.NoError(t, err)
			require.NotNil(t, r)
			require.False(t, r.Ok)
			require.Equal(t, MinKuberVersionNotSpecified, r.Reason)
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
			r, err := NotContainCRDs(tc.uri)
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
			r, err := NotContainCRDs(tc.uri)
			require.NoError(t, err)
			require.NotNil(t, r)
			require.False(t, r.Ok)
			require.Equal(t, ChartContainCRDs, r.Reason)
		})
	}
}
