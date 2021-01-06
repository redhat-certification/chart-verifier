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
	"helmcertifier/pkg/helmcertifier/checks"
)

type CertificateBuilder interface {
	AddResult(result checks.Result) CertificateBuilder
	Build() (Certificate, error)
}

type certificateBuilder struct {
	Results []checks.Result
}

func NewCertificateBuilder() CertificateBuilder {
	return &certificateBuilder{
		Results: []checks.Result{},
	}
}

func (r *certificateBuilder) AddResult(result checks.Result) CertificateBuilder {
	r.Results = append(r.Results, result)
	return r
}

func (r *certificateBuilder) Build() (Certificate, error) {
	res := &certificate{Ok: true}

	for _, cr := range r.Results {
		if !cr.Ok {
			res.Ok = false
			break
		}
	}

	return res, nil
}
