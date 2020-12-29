/*
 * Copyright (C) 29/12/2020, 15:35 igors
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

type Result interface {
	IsOk() bool
}

type result struct {
	Ok bool
}

func (r result) IsOk() bool {
	return r.Ok
}

type ResultBuilder interface {
	AddCheckResult(checkResult CheckResult)
	Build() (CertificationResult, error)
}

type resultBuilder struct {
	CheckResults []CheckResult
}

func NewCertificationResultBuilder() ResultBuilder {
	return &resultBuilder{
		CheckResults: []CheckResult{},
	}
}

func (r *resultBuilder) AddCheckResult(checkResult CheckResult) {
	r.CheckResults = append(r.CheckResults, checkResult)
}

func (r *resultBuilder) Build() (CertificationResult, error) {
	res := result{Ok: true}

	for _, cr := range r.CheckResults {
		if !cr.Ok {
			res.Ok = false
			break
		}
	}

	return res, nil
}
