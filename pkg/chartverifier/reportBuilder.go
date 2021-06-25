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
	"crypto"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	"github.com/redhat-certification/chart-verifier/pkg/chartverifier/profiles"
	"io"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/redhat-certification/chart-verifier/pkg/chartverifier/checks"
	helmchart "helm.sh/helm/v3/pkg/chart"
)

type ReportBuilder interface {
	SetToolVersion(name string) ReportBuilder
	SetProfile(vendorType profiles.VendorType, version string) ReportBuilder
	SetChartUri(name string) ReportBuilder
	AddCheck(check checks.Check, result checks.Result) ReportBuilder
	SetChart(chart *helmchart.Chart) ReportBuilder
	SetCertifiedOpenShiftVersion(version string) ReportBuilder
	Build() (*Report, error)
}

type CheckResult struct {
	checks.Result
	Name string
}

type reportBuilder struct {
	Chart      *helmchart.Chart
	Report     Report
	OCPVersion string
}

func NewReportBuilder() ReportBuilder {
	b := reportBuilder{}
	b.Report = newReport()
	return &b
}

func (r *reportBuilder) SetCertifiedOpenShiftVersion(version string) ReportBuilder {
	r.OCPVersion = version
	return r
}

func (r *reportBuilder) SetToolVersion(version string) ReportBuilder {
	r.Report.Metadata.ToolMetadata.Version = version
	return r
}

func (r *reportBuilder) SetProfile(vendorType profiles.VendorType, version string) ReportBuilder {
	r.Report.Metadata.ToolMetadata.Profile.VendorType = string(vendorType)
	r.Report.Metadata.ToolMetadata.Profile.Version = version
	return r
}

func (r *reportBuilder) SetChartUri(uri string) ReportBuilder {
	r.Report.Metadata.ToolMetadata.ChartUri = uri
	return r
}

func (r *reportBuilder) SetChart(chart *helmchart.Chart) ReportBuilder {
	r.Chart = chart
	r.Report.Metadata.ChartData = chart.Metadata
	return r
}

func (r *reportBuilder) AddCheck(check checks.Check, result checks.Result) ReportBuilder {
	checkReport := r.Report.AddCheck(check)
	checkReport.SetResult(result.Ok, result.Reason)
	return r
}

func (r *reportBuilder) Build() (*Report, error) {

	for _, annotation := range profiles.Get().Annotations {
		switch annotation {
		case profiles.DigestAnnotation:
			r.Report.Metadata.ToolMetadata.Digests.Chart = GenerateSha(r.Chart.Raw)
			r.Report.Metadata.ToolMetadata.Digest = r.Report.Metadata.ToolMetadata.Digests.Chart
		case profiles.LastCertifiedTimestampAnnotation:
			r.Report.Metadata.ToolMetadata.LastCertifiedTimestamp = time.Now().Format("2006-01-02T15:04:05.999999-07:00")
		case profiles.OCPVersionAnnotation:
			if len(r.OCPVersion) == 0 {
				r.Report.Metadata.ToolMetadata.CertifiedOpenShiftVersions = "N/A"
			} else {
				r.Report.Metadata.ToolMetadata.CertifiedOpenShiftVersions = r.OCPVersion
			}
		}
	}

	r.Report.Metadata.ToolMetadata.Digests.Package = GetPackageDigest(r.Report.Metadata.ToolMetadata.ChartUri)

	return &r.Report, nil
}

type By func(p1, p2 *helmchart.File) bool

type fileSorter struct {
	files []*helmchart.File
	by    func(p1, p2 *helmchart.File) bool // Closure used in the Less method.
}

func (by By) sort(files []*helmchart.File) {
	fs := &fileSorter{
		files: files,
		by:    by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(fs)
}

// Len is part of sort.Interface.
func (fs *fileSorter) Len() int {
	return len(fs.files)
}

// Swap is part of sort.Interface.
func (fs *fileSorter) Swap(i, j int) {
	fs.files[i], fs.files[j] = fs.files[j], fs.files[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (fs *fileSorter) Less(i, j int) bool {
	return fs.by(fs.files[i], fs.files[j])
}

func GenerateSha(rawFiles []*helmchart.File) string {

	name := func(f1, f2 *helmchart.File) bool {
		return f1.Name < f2.Name
	}

	chartSha := sha256.New()
	sortedFiles := rawFiles
	By(name).sort(sortedFiles)
	for _, chartFile := range sortedFiles {
		chartSha.Write(chartFile.Data)
	}

	return fmt.Sprintf("sha256:%x", chartSha.Sum(nil))
}

func GetPackageDigest(uri string) string {

	url, err := url.Parse(uri)
	if err != nil {
		return ""
	}
	var chartReader io.Reader
	switch url.Scheme {
	case "http", "https":
		var chartGetResponse *http.Response
		chartGetResponse, err = http.Get(url.String())
		if err == nil {
			chartReader = chartGetResponse.Body
		}
	case "file", "":
		if strings.HasSuffix(url.Path, ".tgz") {
			chartReader, _ = os.Open(url.Path)
		}
	default:
		err = errors.Errorf("scheme %q not supported", url.Scheme)
	}
	if err != nil || chartReader == nil {
		return ""
	}
	return getDigest(chartReader)
}

// Digest hashes a reader and returns a SHA256 digest.
func getDigest(in io.Reader) string {
	if in == nil {
		return ""
	}

	hash := crypto.SHA256.New()
	if _, err := io.Copy(hash, in); err != nil {
		return ""
	}
	return hex.EncodeToString(hash.Sum(nil))
}
