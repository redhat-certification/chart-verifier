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

type CheckResult struct {
	Ok bool
}

type certifier struct {
	registry       *Registry
	requiredChecks []string
}

type CertificationCheckNotFoundErr string

func (c CertificationCheckNotFoundErr) Error() string {
	return "certification check not found: " + string(c)
}

func (c *certifier) Certify(uri string) (CertificationResult, error) {
	result := NewCertificationResultBuilder()

	for _, name := range c.requiredChecks {
		if checkFunc, ok := c.registry.GetCheck(name); !ok {
			return nil, CertificationCheckNotFoundErr(name)
		} else {
			r, err := checkFunc(uri)
			if err != nil {
				return nil, err
			}
			result.AddCheckResult(r)
		}
	}

	return result.Build()
}
