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
	"github.com/Masterminds/semver"
	"github.com/helm/chart-testing/v3/pkg/exec"
	"github.com/redhat-certification/chart-verifier/pkg/chartverifier/checks"
	"github.com/redhat-certification/chart-verifier/pkg/tool"
	"github.com/spf13/viper"
	helmcli "helm.sh/helm/v3/pkg/cli"
)

type CheckNotFoundErr string

func (e CheckNotFoundErr) Error() string {
	return "check not found: " + string(e)
}

type CheckErr string

func (e CheckErr) Error() string {
	return "check error: " + string(e)
}

type OpenShiftVersionErr string

func (e OpenShiftVersionErr) Error() string {
	return "Missing OpenShift version. " + string(e) + ". And the 'openshift-version' flag has not set."
}

type OpenShiftSemVerErr string

func (e OpenShiftSemVerErr) Error() string {
	return "OpenShift version is not following SemVer spec. " + string(e)
}

func NewCheckErr(err error) error {
	return CheckErr(err.Error())
}

// Versioner provides OpenShift version
type Versioner interface {
	getVersion(debug bool) (string, error)
}

type certifier struct {
	config           *viper.Viper
	registry         checks.Registry
	requiredChecks   []string
	settings         *helmcli.EnvSettings
	toolVersion      string
	openshiftVersion string
	values           map[string]interface{}
	version          Versioner
}

func (c *certifier) subConfig(name string) *viper.Viper {
	if sub := c.config.Sub(name); sub == nil {
		return viper.New()
	} else {
		return sub
	}
}

type version struct{}

func (ver *version) getVersion(debug bool) (string, error) {

	procExec := exec.NewProcessExecutor(debug)
	oc := tool.NewOc(procExec)

	// oc.GetVersion() returns an error both in case the oc command can't be executed and
	// the value for the OpenShift version key not present.
	return oc.GetVersion()
}

func (c *certifier) Certify(uri string) (*Certificate, error) {

	chrt, _, err := checks.LoadChartFromURI(uri)
	if err != nil {
		return nil, err
	}

	// oc.GetVersion() returns an error both in case the oc command can't be executed and
	// the value for the OpenShift version key not present.
	osVersion, getVersionErr := c.version.getVersion(c.settings.Debug)

	// From this point on, an error is set and osVersion is empty.
	if getVersionErr != nil && c.openshiftVersion != "" {
		osVersion = c.openshiftVersion
	}

	// osVersion is empty only if an error happened and a default value
	// informed by the user hasn't been informed.
	if osVersion == "" {
		return nil, OpenShiftVersionErr(getVersionErr.Error())
	}

	// osVersion is guaranteed to have a value, not yet validated as a
	// semver value.
	if _, err := semver.NewVersion(osVersion); err != nil {
		return nil, OpenShiftSemVerErr(err.Error())
	}
	// osVersion is guaranteed to a valid value from here onwards.

	result := NewCertificateBuilder().
		SetToolVersion(c.toolVersion).
		SetChartUri(uri).
		SetChart(chrt).
		SetCertifiedOpenShiftVersion(osVersion)

	for _, name := range c.requiredChecks {
		check, ok := c.registry.Get(name)
		if !ok {
			return nil, CheckNotFoundErr(name)
		}

		r, checkErr := check.Func(&checks.CheckOptions{
			HelmEnvSettings: c.settings,
			URI:             uri,
			Values:          c.values,
			ViperConfig:     c.subConfig(name),
		})

		if checkErr != nil {
			return nil, NewCheckErr(checkErr)
		}
		_ = result.AddCheck(name, check.Type, r)

	}

	return result.Build()
}
