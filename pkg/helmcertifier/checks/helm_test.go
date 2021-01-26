/*
 * Copyright (C) 08/01/2021, 01:52, igors
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
	"context"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func serveCharts(ctx context.Context, addr string) {

	mux := http.NewServeMux()
	prefix := "/charts/"
	chartHandler := http.StripPrefix(prefix, http.FileServer(http.Dir("./")))
	mux.Handle(prefix, chartHandler)

	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("server shutdown failed: %s\n", err)
	}
}

func TestLoadChartFromURI(t *testing.T) {
	addr := "127.0.0.1:9876"
	ctx, cancel := context.WithCancel(context.Background())

	type testCase struct {
		description string
		uri         string
	}

	positiveCases := []testCase{
		{
			uri:         "chart-0.1.0-v3.valid.tgz",
			description: "absolute path",
		},
		{
			uri:         "http://" + addr + "/charts/chart-0.1.0-v3.valid.tgz",
			description: "remote path, http",
		},
	}

	go serveCharts(ctx, addr)

	for _, tc := range positiveCases {
		t.Run(tc.description, func(t *testing.T) {
			c, err := loadChartFromURI(tc.uri)
			require.NoError(t, err)
			require.NotNil(t, c)
		})
	}

	cancel()
}
