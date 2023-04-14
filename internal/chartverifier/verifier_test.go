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

	"github.com/redhat-certification/chart-verifier/internal/chartverifier/profiles"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"helm.sh/helm/v3/pkg/cli"

	"github.com/redhat-certification/chart-verifier/internal/chartverifier/checks"
	"github.com/redhat-certification/chart-verifier/internal/testutil"
	apiReport "github.com/redhat-certification/chart-verifier/pkg/chartverifier/report"
)

func isOk(c *apiReport.Report) bool {
	outcome := true
	for _, check := range c.Results {
		if !(check.Outcome == apiReport.PassOutcomeType) {
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

	dummyCheck := checks.Check{CheckId: checks.CheckId{Name: "dummy-check"}}

	erroredCheck := func(_ *checks.CheckOptions) (checks.Result, error) {
		return checks.Result{}, errors.New("artificial error")
	}

	negativeCheck := func(_ *checks.CheckOptions) (checks.Result, error) {
		return checks.Result{Ok: false}, nil
	}

	positiveCheck := func(_ *checks.CheckOptions) (checks.Result, error) {
		return checks.Result{Ok: true}, nil
	}

	validChartURI := "http://" + addr + "/charts/chart-0.1.0-v3.valid.tgz"

	t.Run("Should return error if check does not exist", func(t *testing.T) {
		c := &verifier{
			settings:       cli.New(),
			config:         viper.New(),
			profile:        profiles.Get(),
			registry:       checks.NewRegistry(),
			requiredChecks: []checks.Check{dummyCheck},
		}

		r, err := c.Verify(validChartURI)
		require.Error(t, err)
		require.Nil(t, r)
	})

	t.Run("Should return error if check exists and returns error", func(t *testing.T) {
		dummyCheck.Func = erroredCheck
		c := &verifier{
			settings:       cli.New(),
			config:         viper.New(),
			profile:        profiles.Get(),
			registry:       checks.NewRegistry().Add(dummyCheck.CheckId.Name, "v1.0", erroredCheck),
			requiredChecks: []checks.Check{dummyCheck},
		}

		r, err := c.Verify(validChartURI)
		require.Error(t, err)
		require.Nil(t, r)
	})

	t.Run("Result should be negative if check exists and returns negative", func(t *testing.T) {
		dummyCheck.Func = negativeCheck
		c := &verifier{
			settings:         cli.New(),
			config:           viper.New(),
			profile:          profiles.Get(),
			registry:         checks.NewRegistry().Add(dummyCheck.CheckId.Name, "v1.0", negativeCheck),
			requiredChecks:   []checks.Check{dummyCheck},
			openshiftVersion: "4.9",
		}

		r, err := c.Verify(validChartURI)
		require.NoError(t, err)
		require.NotNil(t, r)
		require.False(t, isOk(r))
	})

	t.Run("Result should be positive if check exists and returns positive", func(t *testing.T) {
		dummyCheck.Func = positiveCheck
		c := &verifier{
			settings:       cli.New(),
			config:         viper.New(),
			profile:        profiles.Get(),
			registry:       checks.NewRegistry().Add(dummyCheck.CheckId.Name, "v1.0", positiveCheck),
			requiredChecks: []checks.Check{dummyCheck},
			webCatalogOnly: true,
		}

		r, err := c.Verify(validChartURI)
		require.NoError(t, err)
		require.NotNil(t, r)
		require.True(t, isOk(r))
	})

	t.Run("Result should be negative is provider deliver is set and uri is not a tarball", func(t *testing.T) {
		dummyCheck.Func = positiveCheck
		c := &verifier{
			settings:       cli.New(),
			config:         viper.New(),
			profile:        profiles.Get(),
			registry:       checks.NewRegistry().Add(dummyCheck.CheckId.Name, "v1.0", positiveCheck),
			requiredChecks: []checks.Check{dummyCheck},
			webCatalogOnly: true,
		}

		r, err := c.Verify("./checks/psql-service-0.1.7")
		require.Error(t, err)
		require.Nil(t, r)
	})
	cancel()
}
