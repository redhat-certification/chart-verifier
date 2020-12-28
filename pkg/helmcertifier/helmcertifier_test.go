/*
 * Copyright (C) 28/12/2020, 16:56 igors
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

package helmcertifier

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCertificationBuilder(t *testing.T) {

	t.Run("Should fail when uri is not set", func(t *testing.T) {
		b := NewCertificationBuilder()

		r, err := b.Build()
		require.Error(t, err)
		require.False(t, r.Ok)
	})

	t.Run("Should fail when checks are not set", func(t *testing.T) {
		b := NewCertificationBuilder()

		r, err := b.SetUri("http://www.example.com/chart.tgz").Build()
		require.Error(t, err)
		require.False(t, r.Ok)
	})

	t.Run("Should succeed when uri and checks are set", func(t *testing.T) {
		b := NewCertificationBuilder()

		r, err := b.
			SetUri("http://www.example.com/chart.tgz").
			SetChecks([]string{"a", "b"}).
			Build()

		require.NoError(t, err)
		require.True(t, r.Ok)
	})
}
