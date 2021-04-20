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
	"github.com/spf13/viper"
)

type CertifierBuilder interface {
	SetRegistry(registry checks.Registry) CertifierBuilder
	SetValues(vals map[string]interface{}) CertifierBuilder
	SetChecks(checks []string) CertifierBuilder
	SetConfig(config *viper.Viper) CertifierBuilder
	SetOverrides([]string) CertifierBuilder
	SetToolVersion(string) CertifierBuilder
	Build() (Certifier, error)
}

type Certifier interface {
	Certify(uri string) (*Certificate, error)
}
