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
	"helm.sh/helm/v3/pkg/cli"
)

type VerifierBuilder interface {
	SetRegistry(registry checks.Registry) VerifierBuilder
	SetValues(vals map[string]interface{}) VerifierBuilder
	SetChecks(checks map[checks.CheckName]checks.Check) VerifierBuilder
	SetConfig(config *viper.Viper) VerifierBuilder
	SetOverrides([]string) VerifierBuilder
	SetToolVersion(string) VerifierBuilder
	SetOpenShiftVersion(string) VerifierBuilder
	SetSettings(settings *cli.EnvSettings) VerifierBuilder
	Build() (Vertifier, error)
}

type Vertifier interface {
	Verify(uri string) (*Report, error)
}
