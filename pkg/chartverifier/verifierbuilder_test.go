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

package chartverifier

import (
	"github.com/redhat-certification/chart-verifier/pkg/chartverifier/checks"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCertificationBuilder(t *testing.T) {

	t.Run("Should fail building verifier when requiredChecks are not set", func(t *testing.T) {
		b := NewVerifierBuilder()

		c, err := b.Build()
		require.Error(t, err)
		require.Nil(t, c)
	})

	t.Run("Should build verifier when requiredChecks are set", func(t *testing.T) {
		b := NewVerifierBuilder()

		c, err := b.
			SetChecks([]checks.CheckName{"a", "b"}).
			Build()

		require.NoError(t, err)
		require.NotNil(t, c)
	})
}
