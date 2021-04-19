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
	"context"
	"errors"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"

	"github.com/redhat-certification/chart-verifier/pkg/chartverifier/checks"
	"github.com/redhat-certification/chart-verifier/pkg/testutil"
)

func TestCertifier_Certify(t *testing.T) {

	addr := "127.0.0.1:9876"
	ctx, cancel := context.WithCancel(context.Background())

	require.NoError(t, testutil.ServeCharts(ctx, addr, "./checks/"))

	dummyCheckName := "dummy-check"

	erroredCheck := func(_ *checks.CheckOptions) (checks.Result, error) {
		return checks.Result{}, errors.New("artificial error")
	}

	negativeCheck := func(_ *checks.CheckOptions) (checks.Result, error) {
		return checks.Result{Ok: false}, nil
	}

	positiveCheck := func(_ *checks.CheckOptions) (checks.Result, error) {
		return checks.Result{Ok: true}, nil
	}

	validChartUri := "http://" + addr + "/charts/chart-0.1.0-v3.valid.tgz"

	t.Run("Should return error if check does not exist", func(t *testing.T) {
		c := &certifier{
			config:         viper.New(),
			registry:       checks.NewRegistry(),
			requiredChecks: []string{dummyCheckName},
		}

		r, err := c.Certify(validChartUri)
		require.Error(t, err)
		require.Nil(t, r)
	})

	t.Run("Should return error if check exists and returns error", func(t *testing.T) {
		c := &certifier{
			config:         viper.New(),
			registry:       checks.NewRegistry().Add(dummyCheckName, erroredCheck),
			requiredChecks: []string{dummyCheckName},
		}

		r, err := c.Certify(validChartUri)
		require.Error(t, err)
		require.Nil(t, r)
	})

	t.Run("Result should be negative if check exists and returns negative", func(t *testing.T) {

		c := &certifier{
			config:         viper.New(),
			registry:       checks.NewRegistry().Add(dummyCheckName, negativeCheck),
			requiredChecks: []string{dummyCheckName},
		}

		r, err := c.Certify(validChartUri)
		require.NoError(t, err)
		require.NotNil(t, r)
		require.False(t, r.IsOk())
	})

	t.Run("Result should be positive if check exists and returns positive", func(t *testing.T) {
		c := &certifier{
			config:         viper.New(),
			registry:       checks.NewRegistry().Add(dummyCheckName, positiveCheck),
			requiredChecks: []string{dummyCheckName},
		}

		r, err := c.Certify(validChartUri)
		require.NoError(t, err)
		require.NotNil(t, r)
		require.True(t, r.IsOk())
	})

	cancel()
}
