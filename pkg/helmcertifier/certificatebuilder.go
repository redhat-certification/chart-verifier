/*
 * Copyright (C) 04/01/2021, 06:48, igors
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
	"errors"

	"helmcertifier/pkg/helmcertifier/checks"
)

type CertificateBuilder interface {
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
	ChartName      string
	ChartVersion   string
	CheckResultMap checkResultMap
}

func NewCertificateBuilder() CertificateBuilder {
	return &certificateBuilder{
		CheckResultMap: checkResultMap{},
	}
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

	return newCertificate(r.ChartName, r.ChartVersion, ok, r.CheckResultMap), nil
}
