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

package checks

import (
	"github.com/spf13/viper"
	helmcli "helm.sh/helm/v3/pkg/cli"
)

type Result struct {
	// Ok indicates whether the result was successful or not.
	Ok bool
	// Reason for the result value.  This is a message indicating
	// the reason for the value of Ok became true or false.
	Reason string
}

func NewResult(outcome bool, reason string) Result {
	result := Result{}
	result.Ok = outcome
	result.Reason = reason
	return result
}

func (r *Result) SetResult(outcome bool, reason string) Result {
	r.Ok = outcome
	r.Reason = reason
	return *r
}

func (r *Result) AddResult(outcome bool, reason string) Result {
	r.Ok = r.Ok && outcome
	if len(r.Reason) > 0 {
		r.Reason += "\n"
	}
	r.Reason += reason
	return *r
}

type AnnotationHolder interface {
	SetCertifiedOpenShiftVersion(version string)
	GetCertifiedOpenShiftVersionFlag() string
}

type CheckName string
type CheckType string
type Version string

const (
	HasReadmeName            CheckName = "has-readme"
	IsHelmV3Name             CheckName = "is-helm-v3"
	ContainsTestName         CheckName = "contains-test"
	ContainsValuesName       CheckName = "contains-values"
	ContainsValuesSchemaName CheckName = "contains-values-schema"
	HasKubeversionName       CheckName = "has-kubeversion"
	NotContainsCRDsName      CheckName = "not-contains-crds"
	HelmLintName             CheckName = "helm-lint"
	NotContainCsiObjectsName CheckName = "not-contain-csi-objects"
	ImagesAreCertifiedName   CheckName = "images-are-certified"
	ChartTestingName         CheckName = "chart-testing"
)

const (
	MandatoryCheckType    CheckType = "Mandatory"
	OptionalCheckType     CheckType = "Optional"
	ExperimentalCheckType CheckType = "Experimental"
)

type CheckId struct {
	Name    CheckName
	Version Version
}
type Check struct {
	CheckId CheckId
	Type    CheckType
	Func    CheckFunc
}

// CheckOptions contains options collected from the environment a check can
// consult to modify its behavior.
type CheckOptions struct {
	// URI is the location of the chart to be checked.
	URI string
	// ViperConfig is the configuration collected by Viper.
	ViperConfig *viper.Viper
	// Values contains the values informed by the user through command line options.
	Values map[string]interface{}
	// HelmEnvSettings contains the Helm related environment settings.
	HelmEnvSettings *helmcli.EnvSettings
	// AnnotationHolder provides and API to set the OpenShift Version
	AnnotationHolder AnnotationHolder
}

type CheckFunc func(options *CheckOptions) (Result, error)

type Registry interface {
	Get(id CheckId) (Check, bool)
	Add(name CheckName, version Version, checkFunc CheckFunc) Registry
	AllChecks() DefaultRegistry
}

type DefaultRegistry map[CheckId]Check

func (r *DefaultRegistry) AllChecks() DefaultRegistry {
	return *r
}

func NewRegistry() Registry {
	return &DefaultRegistry{}
}

func (r *DefaultRegistry) Get(id CheckId) (Check, bool) {
	v, ok := (*r)[id]
	return v, ok
}

func (r *DefaultRegistry) Add(name CheckName, version Version, checkFunc CheckFunc) Registry {

	check := Check{CheckId: CheckId{Name: name, Version: version}, Func: checkFunc}
	(*r)[check.CheckId] = check
	return r
}
