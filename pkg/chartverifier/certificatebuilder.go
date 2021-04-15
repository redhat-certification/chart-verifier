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
	"crypto/sha256"
	"fmt"
	"github.com/redhat-certification/chart-verifier/pkg/chartverifier/checks"
	helmchart "helm.sh/helm/v3/pkg/chart"
	"sort"
	"time"
)

type CertificateBuilder interface {
	SetToolVersion(name string) CertificateBuilder
	SetChartUri(name string) CertificateBuilder
	AddCheck(name string, checkType checks.CheckType, result checks.Result) CertificateBuilder
	SetChart(chart *helmchart.Chart) CertificateBuilder
	Build() (*Certificate, error)
}

type CheckResult struct {
	checks.Result
	Name string
}

type certificateBuilder struct {
	Chart       *helmchart.Chart
	Certificate Certificate
}

func NewCertificateBuilder() CertificateBuilder {
	cB := certificateBuilder{}
	cB.Certificate = newCertificate()
	return &cB
}

func (r *certificateBuilder) SetToolVersion(version string) CertificateBuilder {
	r.Certificate.Metadata.ToolMetadata.Version = version
	return r
}

func (r *certificateBuilder) SetChartUri(uri string) CertificateBuilder {
	r.Certificate.Metadata.ToolMetadata.ChartUri = uri
	return r
}

func (r *certificateBuilder) SetChart(chart *helmchart.Chart) CertificateBuilder {
	r.Chart = chart
	r.Certificate.Metadata.ChartData = chart.Metadata
	return r
}

func (r *certificateBuilder) AddCheck(name string, checkType checks.CheckType, result checks.Result) CertificateBuilder {
	checkReport := r.Certificate.AddCheck(name, checkType)
	checkReport.SetResult(result.Ok, result.Reason)
	return r
}

func (r *certificateBuilder) Build() (*Certificate, error) {

	r.Certificate.Metadata.ToolMetadata.Digest = GenerateSha(r.Chart.Raw)

	r.Certificate.Metadata.ToolMetadata.LastCertifiedTime = time.Now().String()

	return &r.Certificate, nil
}

type By func(p1, p2 *helmchart.File) bool

type fileSorter struct {
	files []*helmchart.File
	by    func(p1, p2 *helmchart.File) bool // Closure used in the Less method.
}

func (by By) Sort(files []*helmchart.File) {
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
	By(name).Sort(sortedFiles)
	for _, chartFile := range sortedFiles {
		chartSha.Write(chartFile.Data)
	}

	return fmt.Sprintf("sha256:%x", chartSha.Sum(nil))
}
