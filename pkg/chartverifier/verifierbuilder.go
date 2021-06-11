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
	"fmt"
	"github.com/redhat-certification/chart-verifier/pkg/chartverifier/profiles"
	"strings"

	"helm.sh/helm/v3/pkg/cli"

	"github.com/spf13/viper"

	"github.com/redhat-certification/chart-verifier/pkg/chartverifier/checks"
)

var defaultRegistry checks.Registry
var profileName string

var initError []string

func init() {
	defaultRegistry = checks.NewRegistry()

	profile := profiles.GetProfile()

	for _, check := range profile.Checks {
		checkIndex := checks.CheckVersion{Name: check.Name, Version: check.Version}
		if newCheck, ok := checks.ChecksMap[checkIndex]; ok {
			newCheck.Type = check.Type
			defaultRegistry.Add(newCheck)
		} else {
			initError = append(initError, fmt.Sprintf("%s : failed to find check: %s, version: %s", profile.Name, check.Name, check.Version))
		}
	}
	profileName = profile.Name
}

func DefaultRegistry() checks.Registry {
	return defaultRegistry
}

type verifierBuilder struct {
	checks           []checks.CheckName
	config           *viper.Viper
	overrides        []string
	registry         checks.Registry
	toolVersion      string
	openshiftVersion string
	values           map[string]interface{}
	settings         *cli.EnvSettings
}

func (b *verifierBuilder) SetSettings(settings *cli.EnvSettings) VerifierBuilder {
	b.settings = settings
	return b
}

func (b *verifierBuilder) SetValues(vals map[string]interface{}) VerifierBuilder {
	b.values = vals
	return b
}

func (b *verifierBuilder) SetRegistry(registry checks.Registry) VerifierBuilder {
	b.registry = registry
	return b
}

func (b *verifierBuilder) SetChecks(checks []checks.CheckName) VerifierBuilder {
	b.checks = checks
	return b
}

func (b *verifierBuilder) SetConfig(config *viper.Viper) VerifierBuilder {
	b.config = config
	return b
}

func (b *verifierBuilder) SetOverrides(overrides []string) VerifierBuilder {
	b.overrides = overrides
	return b
}

func (b *verifierBuilder) SetToolVersion(version string) VerifierBuilder {
	b.toolVersion = version
	return b
}

func (b *verifierBuilder) SetOpenShiftVersion(version string) VerifierBuilder {
	b.openshiftVersion = version
	return b
}

func (b *verifierBuilder) Build() (Vertifier, error) {
	if len(b.checks) == 0 {
		return nil, errors.New("no checks have been required")
	}

	if len(initError) > 0 {
		return nil, errors.New(fmt.Sprintf("Error processing profile %s : %v", profileName, initError))
	}

	if b.registry == nil {
		b.registry = defaultRegistry
	}

	if b.config == nil {
		b.config = viper.New()
	}

	if b.settings == nil {
		b.settings = cli.New()
	}

	// naively override values from the configuration
	for _, val := range b.overrides {
		parts := strings.Split(val, "=")
		b.config.Set(parts[0], parts[1])
	}
	return &verifier{
		config:           b.config,
		registry:         b.registry,
		requiredChecks:   b.checks,
		settings:         b.settings,
		toolVersion:      b.toolVersion,
		profileName:      profileName,
		openshiftVersion: b.openshiftVersion,
		values:           b.values,
	}, nil
}

func NewVerifierBuilder() VerifierBuilder {
	return &verifierBuilder{}
}
