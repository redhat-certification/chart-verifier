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
	"strings"

	"github.com/spf13/viper"

	"github.com/redhat-certification/chart-verifier/pkg/chartverifier/checks"
)

var defaultRegistry checks.Registry

func init() {
	defaultRegistry = checks.NewRegistry()
	defaultRegistry.Add("has-readme", checks.HasReadme)
	defaultRegistry.Add("is-helm-v3", checks.IsHelmV3)
	defaultRegistry.Add("contains-test", checks.ContainsTest)
	defaultRegistry.Add("contains-values", checks.ContainsValues)
	defaultRegistry.Add("contains-values-schema", checks.ContainsValuesSchema)
	defaultRegistry.Add("has-minkubeversion", checks.HasMinKubeVersion)
	defaultRegistry.Add("not-contains-crds", checks.NotContainCRDs)
	defaultRegistry.Add("helm-lint", checks.HelmLint)
	defaultRegistry.Add("not-contain-csi-objects", checks.NotContainCSIObjects)
	defaultRegistry.Add("images-are-certified", checks.ImagesAreCertified)
}

func DefaultRegistry() checks.Registry {
	return defaultRegistry
}

type certifierBuilder struct {
	checks      []string
	config      *viper.Viper
	overrides   []string
	registry    checks.Registry
	toolVersion string
}

func (b *certifierBuilder) SetRegistry(registry checks.Registry) CertifierBuilder {
	b.registry = registry
	return b
}

func (b *certifierBuilder) SetChecks(checks []string) CertifierBuilder {
	b.checks = checks
	return b
}

func (b *certifierBuilder) SetConfig(config *viper.Viper) CertifierBuilder {
	b.config = config
	return b
}

func (b *certifierBuilder) SetOverrides(overrides []string) CertifierBuilder {
	b.overrides = overrides
	return b
}

func (b *certifierBuilder) SetToolVersion(version string) CertifierBuilder {
	b.toolVersion = version
	return b
}

func (b *certifierBuilder) Build() (Certifier, error) {
	if len(b.checks) == 0 {
		return nil, errors.New("no checks have been required")
	}

	if b.registry == nil {
		b.registry = defaultRegistry
	}

	if b.config == nil {
		b.config = viper.New()
	}

	// naively override values from the configuration
	for _, val := range b.overrides {
		parts := strings.Split(val, "=")
		b.config.Set(parts[0], parts[1])
	}

	return &certifier{
		registry:       b.registry,
		requiredChecks: b.checks,
		config:         b.config,
		toolVersion:    b.toolVersion,
	}, nil
}

func NewCertifierBuilder() CertifierBuilder {
	return &certifierBuilder{}
}
