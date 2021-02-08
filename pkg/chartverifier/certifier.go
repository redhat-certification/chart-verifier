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
	"github.com/redhat-certification/chart-verifier/pkg/chartverifier/checks"
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

	chrt, _, err := checks.LoadChartFromURI(uri)
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
