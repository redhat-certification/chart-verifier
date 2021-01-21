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
		{description: "valid tarball, absolute path", uri: "testchart-0.1.0.tgz"},
	}

	for _, tc := range positiveTestCases {
		t.Run(tc.description, func(t *testing.T) {
			r, err := IsHelmV3(tc.uri)
			require.NoError(t, err)
			require.NotNil(t, r)
			require.Equal(t, true, r.Ok)
		})
	}

	errorTestCases := []testCase{
		{description: "invalid tarball, absolute path", uri: "/tmp/chart-v2.tgz"},
		{description: "invalid tarball, relative path", uri: "./chart-v2.tgz"},
		{description: "invalid tarball, http", uri: "http://www.example.com/chart-v2.tgz"},
		{description: "invalid directory, absolute path", uri: "/tmp/chart-v2"},
		{description: "invalid directory, relative path", uri: "./chart-v2"},
	}

	for _, tc := range errorTestCases {
		t.Run(tc.description, func(t *testing.T) {
			r, err := IsHelmV3(tc.uri)
			require.Error(t, err)
			require.NotNil(t, r)
		})
	}

	negativeTestCases := []testCase{
		{description: "invalid tarball, relative path", uri: "helm2testchart-0.1.0.tgz"},
	}

	for _, tc := range negativeTestCases {
		t.Run(tc.description, func(t *testing.T) {
			r, err := IsHelmV3(tc.uri)
			require.NoError(t, err)
			require.NotNil(t, r)
			require.Equal(t, false, r.Ok)
		})
	}
}

func TestHasReadme(t *testing.T) {
	type testCase struct {
		description string
		uri         string
	}

	positiveTestCases := []testCase{
		{description: "tarball contains README.md, absolute path", uri: "/tmp/chart-v3.tgz"},
		{description: "tarball contains README.md, relative path", uri: "./chart-v3.tgz"},
		{description: "tarball contains README.md, http", uri: "http://www.example.com/chart-v3"},
		{description: "directory contains README.md, absolute path", uri: "/tmp/chart-v3"},
		{description: "directory contains README.md, relative path", uri: "./chart-v3"},
	}

	for _, tc := range positiveTestCases {
		t.Run(tc.description, func(t *testing.T) {
			r, err := IsHelmV3(tc.uri)
			require.NoError(t, err)
			require.NotNil(t, r)
		})
	}

	negativeTestCases := []testCase{
		{description: "invalid tarball, absolute path", uri: "/tmp/chart-v3-no-readme.tgz"},
		{description: "invalid tarball, relative path", uri: "./chart-v3-no-readme.tgz"},
		{description: "invalid tarball, http", uri: "http://www.example.com/chart-v3-no-readme.tgz"},
		{description: "invalid directory, absolute path", uri: "/tmp/chart-v3-no-readme"},
		{description: "invalid directory, relative path", uri: "./chart-v3-no-readme"},
	}

	for _, tc := range negativeTestCases {
		t.Run(tc.description, func(t *testing.T) {
			r, err := IsHelmV3(tc.uri)
			require.Error(t, err)
			require.Nil(t, r)
		})
	}
}
