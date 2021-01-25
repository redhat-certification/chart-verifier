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
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"

	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
)

type HTTPChartLoader struct {
	URL *url.URL
}

func (r HTTPChartLoader) Load() (*chart.Chart, error) {
	resp, err := http.Get(r.URL.String())
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return nil, errors.New(string(b))
}

func NewHTTPChartLoader(url *url.URL) loader.ChartLoader {
	return &HTTPChartLoader{URL: url}
}

func loadChartFromURI(uri string) (*chart.Chart, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	switch u.Scheme {
	case "http":
		return NewHTTPChartLoader(u).Load()
	default:
		// try to parse an `uri` string into a `URL` type to decide which `loader.Load*` function to use.
		return loader.LoadFile(u.Path)
	}
}
