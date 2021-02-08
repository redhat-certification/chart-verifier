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
