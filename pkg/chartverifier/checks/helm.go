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

package checks

import (
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"

	"helm.sh/helm/v3/pkg/chartutil"

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

type ChartCache interface {
	MakeKey(uri string) string
	Add(uri string, chrt *chart.Chart) (ChartCacheItem, error)
	Get(uri string) (ChartCacheItem, bool, error)
}

type ChartCacheItem struct {
	Chart *chart.Chart
	Path  string
}

type chartCache struct {
	chartMap map[string]ChartCacheItem
}

func newChartCache() *chartCache {
	return &chartCache{
		chartMap: make(map[string]ChartCacheItem),
	}
}

func (c *chartCache) MakeKey(uri string) string {
	return regexp.MustCompile("[:/?.-]").ReplaceAllString(uri, "_")
}

func (c *chartCache) Get(uri string) (ChartCacheItem, bool, error) {
	if item, ok := c.chartMap[c.MakeKey(uri)]; !ok {
		return ChartCacheItem{}, false, nil
	} else {
		return item, true, nil
	}
}

func (c *chartCache) Add(uri string, chrt *chart.Chart) (ChartCacheItem, error) {
	userCacheDir, err := os.UserCacheDir()
	if err != nil {
		return ChartCacheItem{}, err
	}
	key := c.MakeKey(uri)
	cacheDir := path.Join(userCacheDir, "chart-verifier")
	chartCacheDir := path.Join(cacheDir, key)
	cacheItem := ChartCacheItem{Chart: chrt, Path: chartCacheDir}
	if err = chartutil.SaveDir(chrt, chartCacheDir); err != nil {
		return ChartCacheItem{}, err
	}
	c.chartMap[key] = cacheItem
	return cacheItem, nil
}

var defaultChartCache *chartCache

func init() {
	defaultChartCache = newChartCache()
}

// LoadChartFromURI attempts to retrieve a chart from the given uri string. It accepts "http", "https", "file" schemes,
// and defaults to "file" if there isn't one.
func LoadChartFromURI(uri string) (*chart.Chart, string, error) {
	var (
		chrt *chart.Chart
		err  error
	)

	if cached, ok, _ := defaultChartCache.Get(uri); ok {
		return cached.Chart, cached.Path, nil
	}

	u, err := url.Parse(uri)
	if err != nil {
		return nil, "", err
	}

	switch u.Scheme {
	case "http", "https":
		chrt, err = loadChartFromRemote(u)
	case "file", "":
		chrt, err = loadChartFromAbsPath(u.Path)
	default:
		return nil, "", errors.Errorf("scheme %q not supported", u.Scheme)
	}

	if err != nil {
		return nil, "", err
	}

	if cached, err := defaultChartCache.Add(uri, chrt); err != nil {
		return nil, "", err
	} else {
		return cached.Chart, cached.Path, nil
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
