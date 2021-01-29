/*
 * Copyright (C) 29/12/2020, 15:01 igors
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

type CheckNotFoundErr string

func (e CheckNotFoundErr) Error() string {
	return "check not found: " + string(e)
}

type CheckErr string

func (e CheckErr) Error() string {
	return "check error: " + string(e)
}

func NewCheckErr(err error) error {
	return CheckErr(err.Error())
}

type certifier struct {
	registry       checks.Registry
	requiredChecks []string
}

func (c *certifier) Certify(uri string) (Certificate, error) {

	chrt, err := checks.LoadChartFromURI(uri)
	if err != nil {
		return nil, err
	}

	result := NewCertificateBuilder().
		SetChartName(chrt.Name()).
		SetChartVersion(chrt.AppVersion())

	for _, name := range c.requiredChecks {
		if checkFunc, ok := c.registry.Get(name); !ok {
			return nil, CheckNotFoundErr(name)
		} else {
			r, err := checkFunc(uri)
			if err != nil {
				return nil, NewCheckErr(err)
			}
			_ = result.AddCheckResult(name, r)
		}
	}

	return result.Build()
}
