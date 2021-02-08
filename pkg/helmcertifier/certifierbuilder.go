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

package helmcertifier

import (
	"errors"

	"helmcertifier/pkg/helmcertifier/checks"
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
}

func DefaultRegistry() checks.Registry {
	return defaultRegistry
}

type certifierBuilder struct {
	registry checks.Registry
	checks   []string
}

func (b *certifierBuilder) SetRegistry(registry checks.Registry) CertifierBuilder {
	b.registry = registry
	return b
}

func (b *certifierBuilder) SetChecks(checks []string) CertifierBuilder {
	b.checks = checks
	return b
}

func (b *certifierBuilder) Build() (Certifier, error) {
	if len(b.checks) == 0 {
		return nil, errors.New("no checks have been required")
	}

	if b.registry == nil {
		b.registry = defaultRegistry
	}

	return &certifier{
		registry:       b.registry,
		requiredChecks: b.checks,
	}, nil
}

func NewCertifierBuilder() CertifierBuilder {
	return &certifierBuilder{}
}
