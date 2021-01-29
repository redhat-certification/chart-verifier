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
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
)

// loadChartFromRemote attempts to retrieve a Helm chart from the given remote url. Returns an error if the given url
// doesn't contain the 'http' or 'https' schema, or any other error related to retrieving the contents of the chart.
func loadChartFromRemote(url *url.URL) (*chart.Chart, error) {
	if url.Scheme != "http" && url.Scheme != "https" {
		return nil, errors.Errorf("only 'http' and 'https' schemes are supported, but got %q", url.Scheme)
	}

	resp, err := http.Get(url.String())
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ChartNotFoundErr(url.String())
	}

	return loader.LoadArchive(resp.Body)
}

// loadChartFromAbsPath attempts to retrieve a local Helm chart by resolving the maybe relative path into an absolute
// path from the current working directory.
func loadChartFromAbsPath(path string) (*chart.Chart, error) {
	// although filepath.Abs() can return an error according to its signature, this won't happen (as of go 1.15)
	// because the only invalid value it would accept is an empty string, which is internally converted into "."
	// regardless, the error is still being caught and propagated to avoid being bitten by internal changes in the
	// future
	chartPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	c, err := loader.Load(chartPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ChartNotFoundErr(path)
		}
		return nil, err
	}

	return c, nil
}

// LoadChartFromURI attempts to retrieve a chart from the given uri string. It accepts "http", "https", "file" schemes,
// and defaults to "file" if there isn't one.
func LoadChartFromURI(uri string) (*chart.Chart, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	switch u.Scheme {
	case "http", "https":
		return loadChartFromRemote(u)
	case "file", "":
		return loadChartFromAbsPath(u.Path)
	default:
		return nil, errors.Errorf("scheme %q not supported", u.Scheme)
	}
}

type ChartNotFoundErr string

func (c ChartNotFoundErr) Error() string {
	return "chart not found: " + string(c)
}

func IsChartNotFound(err error) bool {
	_, ok := err.(ChartNotFoundErr)
	return ok
}
