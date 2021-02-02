/*
 * Copyright (C) 29/12/2020, 15:00 igors
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
	"errors"

	"helmcertifier/pkg/helmcertifier/checks"
)

var defaultRegistry checks.Registry

func init() {
	defaultRegistry = checks.NewRegistry()
	defaultRegistry.Add("is-helm-v3", checks.IsHelmV3)
	defaultRegistry.Add("contains-test", checks.ContainsTest)
	defaultRegistry.Add("contains-values", checks.ContainsValues)
	defaultRegistry.Add("contains-values-schema", checks.ContainsValuesSchema)
	defaultRegistry.Add("has-minkubeversion", checks.HasMinKubeVersion)
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
