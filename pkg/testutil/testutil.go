/*
 * Copyright (C) 27/01/2021, 12:51, igors
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

package testutil

import (
	"context"
	"log"
	"net/http"
	"time"
)

// ServeCharts attempts to create a simple HTTP server on the given addr.
func ServeCharts(ctx context.Context, addr string, path string) {
	if path == "" {
		path = "./"
	}

	mux := http.NewServeMux()
	prefix := "/charts/"
	chartHandler := http.StripPrefix(prefix, http.FileServer(http.Dir(path)))
	mux.Handle(prefix, chartHandler)

	srv := &http.Server{Addr: addr, Handler: mux}

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
