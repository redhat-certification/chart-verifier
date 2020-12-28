/*
 * Copyright (C) 28/12/2020, 16:40 igors
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

import "errors"

type CertificationBuilder interface {
	SetUri(uri string) CertificationBuilder
	SetChecks(checks []string) CertificationBuilder
	Build() (Result, error)
}

type Result struct {
	Ok bool
}

type builder struct {
	uri    string
	checks []string
}

func (b *builder) SetUri(uri string) CertificationBuilder {
	b.uri = uri
	return b
}

func (b *builder) SetChecks(checks []string) CertificationBuilder {
	b.checks = checks
	return b
}

func (b *builder) Build() (Result, error) {
	if b.uri == "" {
		return Result{}, errors.New("SetUri() must be called before Build()")
	}
	if len(b.checks) == 0 {
		return Result{}, errors.New("SetChecks() must be called before Build()")
	}

	return Result{Ok: true}, nil
}

func NewCertificationBuilder() CertificationBuilder {
	return &builder{}
}
