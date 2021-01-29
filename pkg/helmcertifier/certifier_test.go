/*
 * Copyright (C) 29/12/2020, 15:13 igors
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
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"helmcertifier/pkg/helmcertifier/checks"
	"helmcertifier/pkg/testutil"
)

func TestCertifier_Certify(t *testing.T) {

	addr := "127.0.0.1:9876"
	ctx, cancel := context.WithCancel(context.Background())
	go testutil.ServeCharts(ctx, addr, "./checks/")

	dummyCheckName := "dummy-check"

	erroredCheck := func(uri string) (checks.Result, error) {
		return checks.Result{}, errors.New("artificial error")
	}

	negativeCheck := func(uri string) (checks.Result, error) {
		return checks.Result{Ok: false}, nil
	}

	positiveCheck := func(uri string) (checks.Result, error) {
		return checks.Result{Ok: true}, nil
	}

	validChartUri := "http://" + addr + "/charts/chart-0.1.0-v3.valid.tgz"

	t.Run("Should return error if check does not exist", func(t *testing.T) {
		c := &certifier{
			registry:       checks.NewRegistry(),
			requiredChecks: []string{dummyCheckName},
		}

		r, err := c.Certify(validChartUri)
		require.Error(t, err)
		require.Nil(t, r)
	})

	t.Run("Should return error if check exists and returns error", func(t *testing.T) {
		c := &certifier{
			registry:       checks.NewRegistry().Add(dummyCheckName, erroredCheck),
			requiredChecks: []string{dummyCheckName},
		}

		r, err := c.Certify(validChartUri)
		require.Error(t, err)
		require.Nil(t, r)
	})

	t.Run("Result should be negative if check exists and returns negative", func(t *testing.T) {

		c := &certifier{
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
