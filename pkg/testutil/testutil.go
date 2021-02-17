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

package testutil

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

// ServeCharts attempts to create a simple HTTP server on the given addr.
func ServeCharts(ctx context.Context, addr string, path string) error {
	if path == "" {
		path = "./"
	}

	mux := http.NewServeMux()
	prefix := "/charts/"
	chartHandler := http.StripPrefix(prefix, http.FileServer(http.Dir(path)))
	mux.Handle(prefix, chartHandler)

	srv := &http.Server{Addr: addr, Handler: mux}

	// listen and server are separated here to catch listen issues before serve executes, otherwise listen errors can't
	// be propagated to the caller.
	ln, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}

	go func() {
		if err := srv.Serve(ln); err != nil && err != http.ErrServerClosed {
			log.Fatalf("serve: %s\n", err)
		}
	}()

	// spawn an extra gofunc to shutdown the server once the context is cancelled
	go func() {
		<-ctx.Done()

		ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer func() {
			cancel()
		}()

		if err := srv.Shutdown(ctxShutdown); err != nil {
			log.Fatalf("server shutdown failed: %s\n", err)
		}
	}()

	return nil
}
