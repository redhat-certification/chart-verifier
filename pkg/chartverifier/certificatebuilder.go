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
	"errors"

	"github.com/redhat-certification/chart-verifier/pkg/chartverifier/checks"
)

type CertificateBuilder interface {
	SetToolVersion(name string) CertificateBuilder
	SetChartUri(name string) CertificateBuilder
	SetChartName(name string) CertificateBuilder
	SetChartVersion(version string) CertificateBuilder
	AddCheckResult(name string, result checks.Result) CertificateBuilder
	Build() (Certificate, error)
}

type CheckResult struct {
	checks.Result
	Name string
}

type certificateBuilder struct {
	ToolVersion    string
	ChartUri       string
	ChartName      string
	ChartVersion   string
	CheckResultMap checkResultMap
}

func NewCertificateBuilder() CertificateBuilder {
	return &certificateBuilder{
		CheckResultMap: checkResultMap{},
	}
}

func (r *certificateBuilder) SetToolVersion(version string) CertificateBuilder {
	r.ToolVersion = version
	return r
}

func (r *certificateBuilder) SetChartUri(uri string) CertificateBuilder {
	r.ChartUri = uri
	return r
}

func (r *certificateBuilder) SetChartName(name string) CertificateBuilder {
	r.ChartName = name
	return r
}

func (r *certificateBuilder) SetChartVersion(version string) CertificateBuilder {
	r.ChartVersion = version
	return r
}

func (r *certificateBuilder) AddCheckResult(name string, result checks.Result) CertificateBuilder {
	r.CheckResultMap[name] = checkResult{Ok: result.Ok, Reason: result.Reason}
	return r
}

func (r *certificateBuilder) Build() (Certificate, error) {

	if r.ChartName == "" {
		return nil, errors.New("chart name must be set")
	}

	if r.ChartVersion == "" {
		return nil, errors.New("chart version must be set")
	}

	ok := true

	for _, v := range r.CheckResultMap {
		if !v.Ok {
			ok = false
			break
		}
	}

	return newCertificate(r.ChartName, r.ChartVersion, r.ChartUri, r.ToolVersion, ok, r.CheckResultMap), nil
}
