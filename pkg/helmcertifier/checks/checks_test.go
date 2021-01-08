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
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIsHelmV3(t *testing.T) {
	type testCase struct {
		description string
		uri         string
	}

	testCases := []testCase{
		{description: "absolute path", uri: "/tmp/chart.tgz"},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			r, err := IsHelmV3(tc.uri)
			require.NoError(t, err)
			require.NotNil(t, r)
		})
	}
}
