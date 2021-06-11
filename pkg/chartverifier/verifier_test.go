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
	"helm.sh/helm/v3/pkg/cli"

	"github.com/redhat-certification/chart-verifier/pkg/chartverifier/checks"
	"github.com/redhat-certification/chart-verifier/pkg/testutil"
)

func (c *Report) isOk() bool {
	outcome := true
	for _, check := range c.Results {
		if !(check.Outcome == PassOutcomeType) {
			outcome = false
			break
		}
	}
	return outcome
}

func TestVerifier_Verify(t *testing.T) {

	addr := "127.0.0.1:9876"
	ctx, cancel := context.WithCancel(context.Background())

	require.NoError(t, testutil.ServeCharts(ctx, addr, "./checks/"))

	var dummyCheckName checks.CheckName = "dummy-check"

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
		c := &verifier{
			settings:       cli.New(),
			config:         viper.New(),
			registry:       checks.NewRegistry(),
			requiredChecks: []checks.CheckName{dummyCheckName},
		}

		r, err := c.Verify(validChartUri)
		require.Error(t, err)
		require.Nil(t, r)
	})

	t.Run("Should return error if check exists and returns error", func(t *testing.T) {
		c := &verifier{
			settings:       cli.New(),
			config:         viper.New(),
			registry:       checks.NewRegistry().Add(checks.Check{Name: dummyCheckName, Type: checks.MandatoryCheckType, Func: erroredCheck}),
			requiredChecks: []checks.CheckName{dummyCheckName},
		}

		r, err := c.Verify(validChartUri)
		require.Error(t, err)
		require.Nil(t, r)
	})

	t.Run("Result should be negative if check exists and returns negative", func(t *testing.T) {

		c := &verifier{
			settings:         cli.New(),
			config:           viper.New(),
			registry:         checks.NewRegistry().Add(checks.Check{Name: dummyCheckName, Type: checks.MandatoryCheckType, Func: negativeCheck}),
			requiredChecks:   []checks.CheckName{dummyCheckName},
			openshiftVersion: "4.9",
		}

		r, err := c.Verify(validChartUri)
		require.NoError(t, err)
		require.NotNil(t, r)
		require.False(t, r.isOk())
	})

	t.Run("Result should be positive if check exists and returns positive", func(t *testing.T) {
		c := &verifier{
			settings:       cli.New(),
			config:         viper.New(),
			registry:       checks.NewRegistry().Add(checks.Check{Name: dummyCheckName, Type: checks.MandatoryCheckType, Func: positiveCheck}),
			requiredChecks: []checks.CheckName{dummyCheckName},
		}

		r, err := c.Verify(validChartUri)
		require.NoError(t, err)
		require.NotNil(t, r)
		require.True(t, r.isOk())
	})

	cancel()
}
