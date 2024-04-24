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

package cmd

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/redhat-certification/chart-verifier/internal/chartverifier/checks"
	"github.com/redhat-certification/chart-verifier/internal/chartverifier/utils"
	apiChecks "github.com/redhat-certification/chart-verifier/pkg/chartverifier/checks"
	apiReport "github.com/redhat-certification/chart-verifier/pkg/chartverifier/report"
)

func TestCertify(t *testing.T) {
	tests := []struct {
		name               string
		args               []string
		validateErrorFunc  func(error)
		validateOutputFunc func(*bytes.Buffer)
	}{
		{
			name: "Should fail when no argument is given",
			args: nil,
			validateErrorFunc: func(err error) {
				require.Error(t, err)
			},
		},
		{
			name: "Should fail when chart does not exist and argument is given",
			args: []string{
				"-E",
				"../internal/chartverifier/checks/chart-0.1.0-v3.non-existing.tgz",
			},
			validateErrorFunc: func(err error) {
				require.Error(t, err)
				require.True(t, checks.IsChartNotFound(err))
			},
		},
		{
			name: "Should fail when the chart does not exist for empty set of checks",
			args: []string{
				"-E",
				"../internal/chartverifier/checks/chart-0.1.0-v3.non-existing.tgz",
				"-o",
			},
			validateErrorFunc: func(err error) {
				require.Error(t, err)
				require.False(t, checks.IsChartNotFound(err))
			},
		},
		{
			name: "Should fail when the chart does not exist for single check",
			args: []string{
				"-e", "is-helm-vv3",
				"-E",
				"../internal/chartverifier/checks/chart-0.1.0-v3.non-existing.tgz",
			},
			validateErrorFunc: func(err error) {
				require.Error(t, err)
				require.False(t, checks.IsChartNotFound(err))
			},
		},
		{
			name: "Should fail when the chart exists but the single check does not",
			args: []string{
				"-e", "is-helm-vv3",
				"-E",
				"../internal/chartverifier/checks/chart-0.1.0-v3.valid.tgz",
			},
			validateErrorFunc: func(err error) {
				require.Error(t, err)
				require.False(t, checks.IsChartNotFound(err))
			},
		},
		{
			name: "Should succeed when the chart exists and is valid for a single check",
			args: []string{
				"-e", "is-helm-v3",
				"-V", "4.9",
				"-E",
				"../internal/chartverifier/checks/chart-0.1.0-v3.valid.tgz",
			},
			validateErrorFunc: func(err error) {
				require.NoError(t, err)
			},
			validateOutputFunc: func(output *bytes.Buffer) {
				require.NotEmpty(t, output.String())

				expected := "results:\n" +
					"    - check: v1.0/is-helm-v3\n" +
					"      type: Mandatory\n" +
					"      outcome: PASS\n" +
					"      reason: API version is V2, used in Helm 3\n"
				require.Contains(t, output.String(), expected)
			},
		},
		{
			name: "Should succeed when the chart exists and a chart value is overridden",
			args: []string{
				"-e", "helm-lint",
				"-S", "replicaCount=2",
				"-E",
				"../internal/chartverifier/checks/chart-0.1.0-v3.with-additionalproperties-false.tgz",
			},
			validateErrorFunc: func(err error) {
				require.NoError(t, err)
			},
			validateOutputFunc: func(output *bytes.Buffer) {
				require.NotEmpty(t, output.String())

				expected := "results:\n" +
					"    - check: v1.0/helm-lint\n" +
					"      type: Mandatory\n" +
					"      outcome: PASS\n" +
					"      reason: Helm lint successful\n"
				require.Contains(t, output.String(), expected)
			},
		},
		{
			name: "Should display JSON certificate when option --output and argument values are given",
			args: []string{
				"-e", "is-helm-v3", // only consider a single check, perhaps more checks in the future
				"-V", "4.9",
				"-o", "json",
				"-E",
				"../internal/chartverifier/checks/chart-0.1.0-v3.valid.tgz",
			},
			validateErrorFunc: func(err error) {
				require.NoError(t, err)
			},
			validateOutputFunc: func(output *bytes.Buffer) {
				require.NotEmpty(t, output.String())

				// attempts to deserialize the command's output, expecting a json string
				certificate := apiReport.Report{}

				err := json.Unmarshal(output.Bytes(), &certificate)
				require.NoError(t, err)
				require.True(t, len(certificate.Results) == 1, "Expected only 1 result")
				require.Equal(t, certificate.Results[0].Check, apiChecks.CheckName("v1.0/is-helm-v3"))
				require.Equal(t, certificate.Results[0].Outcome, apiReport.PassOutcomeType)
				require.Equal(t, certificate.Results[0].Type, apiChecks.MandatoryCheckType)
				require.Equal(t, certificate.Results[0].Reason, checks.Helm3Reason)
			},
		},
		{
			name: "Should display YAML certificate when option --output and argument values are given",
			args: []string{
				"-e", "is-helm-v3", // only consider a single check, perhaps more checks in the future
				"-V", "4.9",
				"-o", "yaml",
				"../internal/chartverifier/checks/chart-0.1.0-v3.valid.tgz",
				"-E",
			},
			validateErrorFunc: func(err error) {
				require.NoError(t, err)
			},
			validateOutputFunc: func(output *bytes.Buffer) {
				require.NotEmpty(t, output.String())

				// attempts to deserialize the command's output, expecting a json string
				certificate := apiReport.Report{}
				err := yaml.Unmarshal(output.Bytes(), &certificate)
				require.NoError(t, err)
				require.True(t, len(certificate.Results) == 1, "Expected only 1 result")
				require.Equal(t, certificate.Results[0].Check, apiChecks.CheckName("v1.0/is-helm-v3"))
				require.Equal(t, certificate.Results[0].Outcome, apiReport.PassOutcomeType)
				require.Equal(t, certificate.Results[0].Type, apiChecks.MandatoryCheckType)
				require.Equal(t, certificate.Results[0].Reason, checks.Helm3Reason)
			},
		},
		{
			name: "Should see webCatalogOnly is true for -W flag and chart-uri is not set",
			args: []string{
				"-e", "has-readme", // only consider a single check, perhaps more checks in the future
				"../internal/chartverifier/checks/chart-0.1.0-v3.valid.tgz",
				"-W",
				"-E",
			},
			validateErrorFunc: func(err error) {
				require.NoError(t, err)
			},
			validateOutputFunc: func(output *bytes.Buffer) {
				require.NotEmpty(t, output.String())

				// attempts to deserialize the command's output, expecting a json string
				certificate := apiReport.Report{}
				err := yaml.Unmarshal(output.Bytes(), &certificate)
				require.NoError(t, err)
				require.True(t, certificate.Metadata.ToolMetadata.WebCatalogOnly)
				require.True(t, certificate.Metadata.ToolMetadata.ChartUri == "N/A")
			},
		},
		{
			name: "Should see webCatalogOnly is false if no -W flag and chart-uri is set",
			args: []string{
				"-e", "has-readme", // only consider a single check, perhaps more checks in the future
				"../internal/chartverifier/checks/chart-0.1.0-v3.valid.tgz",
				"-E",
			},
			validateErrorFunc: func(err error) {
				require.NoError(t, err)
			},
			validateOutputFunc: func(output *bytes.Buffer) {
				require.NotEmpty(t, output.String())

				// attempts to deserialize the command's output, expecting a json string
				certificate := apiReport.Report{}
				err := yaml.Unmarshal(output.Bytes(), &certificate)
				require.NoError(t, err)
				require.False(t, certificate.Metadata.ToolMetadata.WebCatalogOnly)
				require.True(t, certificate.Metadata.ToolMetadata.ChartUri == "../internal/chartverifier/checks/chart-0.1.0-v3.valid.tgz")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cmd := NewVerifyCmd(viper.New())
			outBuf := bytes.NewBufferString("")
			cmd.SetOut(outBuf)
			errBuf := bytes.NewBufferString("")
			cmd.SetErr(errBuf)
			if tc.validateOutputFunc != nil {
				utils.CmdStdout = outBuf
			}

			if tc.args != nil {
				cmd.SetArgs(tc.args)
			}

			err := cmd.Execute()
			if tc.validateErrorFunc != nil {
				tc.validateErrorFunc(err)
			}

			if tc.validateOutputFunc != nil {
				tc.validateOutputFunc(outBuf)
			}
		})
	}
}

func TestBuildChecks(t *testing.T) {
	defaultCheckFunc := func(enabledSet, disabledSet []apiChecks.CheckName, err error) {
		require.Error(t, err)
		require.Nil(t, enabledSet)
		require.Nil(t, disabledSet)
	}

	tests := []struct {
		name               string
		enabledChecks      []string
		disabledChecks     []string
		validateChecksFunc func([]apiChecks.CheckName, []apiChecks.CheckName, error)
	}{
		{
			name:               "Should fail when enabledChecks and disabledChecks have more than one item at the same time",
			enabledChecks:      []string{string(apiChecks.HasReadme)},
			disabledChecks:     []string{string(apiChecks.ChartTesting)},
			validateChecksFunc: defaultCheckFunc,
		},
		{
			name:               "Should fail when enabled check is unknown",
			enabledChecks:      []string{},
			disabledChecks:     []string{"d"},
			validateChecksFunc: defaultCheckFunc,
		},
		{
			name:               "Should fail when disabled check is unknown",
			enabledChecks:      []string{"e"},
			disabledChecks:     []string{},
			validateChecksFunc: defaultCheckFunc,
		},
		{
			name:           "Should return no checks when neither enabled or disabled checks have been informed",
			enabledChecks:  []string{},
			disabledChecks: []string{},
			validateChecksFunc: func(enabledSet, disabledSet []apiChecks.CheckName, err error) {
				require.NoError(t, err)
				require.Nil(t, enabledSet)
				require.Nil(t, disabledSet)
			},
		},
		{
			name:           "Should return enabled checks",
			enabledChecks:  []string{"has-readme", "has-kubeversion", "images-are-certified"},
			disabledChecks: []string{},
			validateChecksFunc: func(enabledSet, disabledSet []apiChecks.CheckName, err error) {
				require.NoError(t, err)
				require.True(t, len(enabledSet) == 3)
				require.Nil(t, disabledSet)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			enabledSet, disabledSet, err := buildChecks(tc.enabledChecks, tc.disabledChecks)
			tc.validateChecksFunc(enabledSet, disabledSet, err)
		})
	}
}

func TestSignatureCheck(t *testing.T) {
	tests := []struct {
		name               string
		args               []string
		validateErrorFunc  func(error)
		validateOutputFunc func(*bytes.Buffer)
	}{
		{
			name: "Unsigned chart no key should pass",
			args: []string{
				"-e", "signature-is-valid", // only consider a single check, perhaps more checks in the future
				"../internal/chartverifier/checks/chart-0.1.0-v3.valid.tgz",
				"-E",
			},
			validateErrorFunc: func(err error) {
				require.NoError(t, err)
			},
			validateOutputFunc: func(output *bytes.Buffer) {
				require.NotEmpty(t, output.String())

				require.Contains(t, output.String(), "outcome: SKIPPED")
				require.Contains(t, output.String(), checks.ChartNotSigned)
				require.Contains(t, output.String(), checks.SignatureIsNotPresentSuccess)
			},
		},
		{
			name: "Unsigned chart with key should pass",
			args: []string{
				"-e", "signature-is-valid", // only consider a single check, perhaps more checks in the future
				"-k", "../tests/charts/psql-service/0.1.11/psql-service-0.1.11.tgz.key",
				"../internal/chartverifier/checks/chart-0.1.0-v3.valid.tgz",
				"-E",
			},
			validateErrorFunc: func(err error) {
				require.NoError(t, err)
			},
			validateOutputFunc: func(output *bytes.Buffer) {
				require.NotEmpty(t, output.String())

				require.Contains(t, output.String(), "outcome: SKIPPED")
				require.Contains(t, output.String(), checks.ChartNotSigned)
				require.Contains(t, output.String(), checks.SignatureIsNotPresentSuccess)
			},
		},
		{
			name: "Signed chart with key should pass",
			args: []string{
				"-e", "signature-is-valid", // only consider a single check, perhaps more checks in the future
				"-k", "../tests/charts/psql-service/0.1.11/psql-service-0.1.11.tgz.key",
				"../tests/charts/psql-service/0.1.11/psql-service-0.1.11.tgz",
				"-E",
			},
			validateErrorFunc: func(err error) {
				require.NoError(t, err)
			},
			validateOutputFunc: func(output *bytes.Buffer) {
				require.NotEmpty(t, output.String())

				require.Contains(t, output.String(), "outcome: PASS")
				require.Contains(t, output.String(), checks.ChartSigned)
				require.Contains(t, output.String(), checks.SignatureIsValidSuccess)
			},
		},
		{
			name: "Signed chart with bad key should fail",
			args: []string{
				"-e", "signature-is-valid", // only consider a single check, perhaps more checks in the future
				"-k", "../tests/charts/psql-service/0.1.11/psql-service-0.1.11.tgz.badkey",
				"../tests/charts/psql-service/0.1.11/psql-service-0.1.11.tgz",
				"-E",
			},
			validateErrorFunc: func(err error) {
				require.NoError(t, err)
			},
			validateOutputFunc: func(output *bytes.Buffer) {
				require.NotEmpty(t, output.String())

				require.Contains(t, output.String(), "outcome: FAIL")
				require.Contains(t, output.String(), checks.ChartSigned)
				require.Contains(t, output.String(), checks.SignatureFailure)
			},
		},
		{
			name: "Signed chart with no key should skip",
			args: []string{
				"-e", "signature-is-valid",
				"https://github.com/redhat-certification/chart-verifier/blob/main/tests/charts/psql-service/0.1.11/psql-service-0.1.11.tgz?raw=true",
				"-E",
			},
			validateErrorFunc: func(err error) {
				require.NoError(t, err)
			},
			validateOutputFunc: func(output *bytes.Buffer) {
				require.NotEmpty(t, output.String())

				require.Contains(t, output.String(), "outcome: SKIPPED")
				require.Contains(t, output.String(), checks.ChartSigned)
				require.Contains(t, output.String(), checks.SignatureNoKey)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cmd := NewVerifyCmd(viper.New())
			outBuf := bytes.NewBufferString("")
			if tc.validateOutputFunc != nil {
				utils.CmdStdout = outBuf
			}
			errBuf := bytes.NewBufferString("")
			cmd.SetErr(errBuf)

			if tc.args != nil {
				cmd.SetArgs(tc.args)
			}

			err := cmd.Execute()
			if tc.validateErrorFunc != nil {
				tc.validateErrorFunc(err)
			}

			if tc.validateOutputFunc != nil {
				tc.validateOutputFunc(outBuf)
			}
		})
	}
}
