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

type CheckType string

const (
	MandatoryCheckType    CheckType = "Mandatory"
	OptionalCheckType     CheckType = "Optional"
	ExperimentalCheckType CheckType = "Experimental"
)

type Check struct {
	Name CheckName
	Type CheckType
	Func CheckFunc
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
	Get(name CheckName) (Check, bool)
	Add(check Check) Registry
	AllChecks() []CheckName
}

type defaultRegistry map[CheckName]Check

type CheckVersion struct {
	Name    CheckName
	Version string
}

var ChecksMap map[CheckVersion]Check

func init() {

	ChecksMap = make(map[CheckVersion]Check)

	ChecksMap[CheckVersion{Name: HasReadmeName, Version: "1.0"}] = Check{Name: HasReadmeName, Func: HasReadme}
	ChecksMap[CheckVersion{Name: IsHelmV3Name, Version: "1.0"}] = Check{Name: IsHelmV3Name, Func: IsHelmV3}
	ChecksMap[CheckVersion{Name: ContainsTestName, Version: "1.0"}] = Check{Name: ContainsTestName, Func: ContainsTest}
	ChecksMap[CheckVersion{Name: ContainsValuesName, Version: "1.0"}] = Check{Name: ContainsValuesName, Func: ContainsValues}
	ChecksMap[CheckVersion{Name: ContainsValuesSchemaName, Version: "1.0"}] = Check{Name: ContainsValuesSchemaName, Func: ContainsValuesSchema}
	ChecksMap[CheckVersion{Name: HasKubeversionName, Version: "1.0"}] = Check{Name: HasKubeversionName, Func: IsHelmV3}
	ChecksMap[CheckVersion{Name: NotContainsCRDsName, Version: "1.0"}] = Check{Name: NotContainsCRDsName, Func: NotContainCRDs}
	ChecksMap[CheckVersion{Name: HelmLintName, Version: "1.0"}] = Check{Name: HelmLintName, Func: HelmLint}
	ChecksMap[CheckVersion{Name: NotContainCsiObjectsName, Version: "1.0"}] = Check{Name: NotContainCsiObjectsName, Func: NotContainCSIObjects}
	ChecksMap[CheckVersion{Name: ImagesAreCertifiedName, Version: "1.0"}] = Check{Name: ImagesAreCertifiedName, Func: ImagesAreCertified}
	ChecksMap[CheckVersion{Name: ChartTestingName, Version: "1.0"}] = Check{Name: ChartTestingName, Func: ChartTesting}
}

func (r *defaultRegistry) AllChecks() []CheckName {
	allChecks := make([]CheckName, 0)
	for k, _ := range *r {
		allChecks = append(allChecks, k)
	}
	return allChecks
}

func NewRegistry() Registry {
	return &defaultRegistry{}
}

func (r *defaultRegistry) Get(name CheckName) (Check, bool) {
	v, ok := (*r)[name]
	return v, ok
}

func (r *defaultRegistry) Add(check Check) Registry {
	(*r)[check.Name] = check
	return r
}
