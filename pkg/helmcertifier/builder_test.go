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

	t.Run("Should fail building certifier when requiredChecks are not set", func(t *testing.T) {
		b := NewCertifierBuilder()

		c, err := b.Build()
		require.Error(t, err)
		require.Nil(t, c)
	})

	t.Run("Should build certifier when requiredChecks are set", func(t *testing.T) {
		b := NewCertifierBuilder()

		c, err := b.
			SetChecks([]string{"a", "b"}).
			Build()

		require.NoError(t, err)
		require.NotNil(t, c)
	})
}
