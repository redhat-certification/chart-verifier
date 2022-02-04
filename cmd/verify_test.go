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
	"fmt"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/redhat-certification/chart-verifier/pkg/chartverifier"
	"github.com/redhat-certification/chart-verifier/pkg/chartverifier/checks"
	"github.com/redhat-certification/chart-verifier/pkg/chartverifier/utils"
)

func TestCertify(t *testing.T) {

	t.Run("Should fail when no argument is given", func(t *testing.T) {
		cmd := NewVerifyCmd(viper.New())
		outBuf := bytes.NewBufferString("")
		cmd.SetOut(outBuf)
		errBuf := bytes.NewBufferString("")
		cmd.SetErr(errBuf)

		require.Error(t, cmd.Execute())
	})

	t.Run("Should fail when chart does not exist and argument is given", func(t *testing.T) {
		cmd := NewVerifyCmd(viper.New())
		outBuf := bytes.NewBufferString("")
		cmd.SetOut(outBuf)
		errBuf := bytes.NewBufferString("")
		cmd.SetErr(errBuf)

		cmd.SetArgs([]string{"-E", "../pkg/chartverifier/checks/chart-0.1.0-v3.non-existing.tgz"})

		err := cmd.Execute()
		require.Error(t, err)
		require.True(t, checks.IsChartNotFound(err))
	})

	t.Run("Should fail when the chart does not exist for empty set of checks", func(t *testing.T) {
		cmd := NewVerifyCmd(viper.New())
		outBuf := bytes.NewBufferString("")
		cmd.SetOut(outBuf)
		errBuf := bytes.NewBufferString("")
		cmd.SetErr(errBuf)

		cmd.SetArgs([]string{"-E", "../pkg/chartverifier/checks/chart-0.1.0-v3.non-existing.tgz", "-o"})
		err := cmd.Execute()
		require.Error(t, err)
		require.False(t, checks.IsChartNotFound(err))
	})

	t.Run("Should fail when the chart does not exist for single check", func(t *testing.T) {
		cmd := NewVerifyCmd(viper.New())
		outBuf := bytes.NewBufferString("")
		cmd.SetOut(outBuf)
		errBuf := bytes.NewBufferString("")
		cmd.SetErr(errBuf)

		cmd.SetArgs([]string{
			"-e", "is-helm-vv3",
			"-E",
			"../pkg/chartverifier/checks/chart-0.1.0-v3.non-existing.tgz",
		})
		err := cmd.Execute()
		require.Error(t, err)
		require.False(t, checks.IsChartNotFound(err))
	})

	t.Run("Should fail when the chart exists but the single check does not", func(t *testing.T) {
		cmd := NewVerifyCmd(viper.New())
		outBuf := bytes.NewBufferString("")
		cmd.SetOut(outBuf)
		errBuf := bytes.NewBufferString("")
		cmd.SetErr(errBuf)

		cmd.SetArgs([]string{
			"-e", "is-helm-vv3",
			"-E",
			"../pkg/chartverifier/checks/chart-0.1.0-v3.valid.tgz",
		})
		err := cmd.Execute()
		require.Error(t, err)
		require.False(t, checks.IsChartNotFound(err))
	})

	t.Run("Should succeed when the chart exists and is valid for a single check", func(t *testing.T) {
		cmd := NewVerifyCmd(viper.New())

		errBuf := bytes.NewBufferString("")
		cmd.SetErr(errBuf)
		outBuf := bytes.NewBufferString("")
		utils.CmdStdout = outBuf

		cmd.SetArgs([]string{
			"-e", "is-helm-v3",
			"-V", "4.9",
			"-E",
			"../pkg/chartverifier/checks/chart-0.1.0-v3.valid.tgz",
		})

		require.NoError(t, cmd.Execute())
		require.NotEmpty(t, outBuf.String())

		expected := "results:\n" +
			"    - check: v1.0/is-helm-v3\n" +
			"      type: Mandatory\n" +
			"      outcome: PASS\n" +
			"      reason: API version is V2, used in Helm 3\n"
		require.Contains(t, outBuf.String(), expected)
	})

	t.Run("Should display JSON certificate when option --output and argument values are given", func(t *testing.T) {
		cmd := NewVerifyCmd(viper.New())
		errBuf := bytes.NewBufferString("")
		cmd.SetErr(errBuf)
		outBuf := bytes.NewBufferString("")
		utils.CmdStdout = outBuf

		cmd.SetArgs([]string{
			"-e", "is-helm-v3", // only consider a single check, perhaps more checks in the future
			"-V", "4.9",
			"-o", "json",
			"-E",
			"../pkg/chartverifier/checks/chart-0.1.0-v3.valid.tgz",
		})

		require.NoError(t, cmd.Execute())
		require.NotEmpty(t, outBuf.String())

		// attempts to deserialize the command's output, expecting a json string
		certificate := chartverifier.Report{}

		err := json.Unmarshal([]byte(outBuf.String()), &certificate)
		require.NoError(t, err)
		require.True(t, len(certificate.Results) == 1, "Expected only 1 result")
		require.Equal(t, certificate.Results[0].Check, checks.CheckName("v1.0/is-helm-v3"))
		require.Equal(t, certificate.Results[0].Outcome, chartverifier.PassOutcomeType)
		require.Equal(t, certificate.Results[0].Type, checks.MandatoryCheckType)
		require.Equal(t, certificate.Results[0].Reason, checks.Helm3Reason)
	})

	t.Run("Should display YAML certificate when option --output and argument values are given", func(t *testing.T) {
		cmd := NewVerifyCmd(viper.New())
		outBuf := bytes.NewBufferString("")
		utils.CmdStdout = outBuf
		errBuf := bytes.NewBufferString("")
		cmd.SetErr(errBuf)

		cmd.SetArgs([]string{
			"-e", "is-helm-v3", // only consider a single check, perhaps more checks in the future
			"-V", "4.9",
			"-o", "yaml",
			"../pkg/chartverifier/checks/chart-0.1.0-v3.valid.tgz",
			"-E",
		})
		require.NoError(t, cmd.Execute())
		require.NotEmpty(t, outBuf.String())

		// attempts to deserialize the command's output, expecting a json string
		certificate := chartverifier.Report{}
		err := yaml.Unmarshal([]byte(outBuf.String()), &certificate)
		require.NoError(t, err)
		require.True(t, len(certificate.Results) == 1, "Expected only 1 result")
		require.Equal(t, certificate.Results[0].Check, checks.CheckName("v1.0/is-helm-v3"))
		require.Equal(t, certificate.Results[0].Outcome, chartverifier.PassOutcomeType)
		require.Equal(t, certificate.Results[0].Type, checks.MandatoryCheckType)
		require.Equal(t, certificate.Results[0].Reason, checks.Helm3Reason)

	})

	t.Run("should see providerControlledDelivery is true for -d flag", func(t *testing.T) {

		cmd := NewVerifyCmd(viper.New())
		outBuf := bytes.NewBufferString("")
		utils.CmdStdout = outBuf
		errBuf := bytes.NewBufferString("")
		cmd.SetErr(errBuf)

		cmd.SetArgs([]string{
			"-e", "has-readme", // only consider a single check, perhaps more checks in the future
			"../pkg/chartverifier/checks/chart-0.1.0-v3.valid.tgz",
			"-d",
			"-E"})

		require.NoError(t, cmd.Execute())
		require.NotEmpty(t, outBuf.String())

		// attempts to deserialize the command's output, expecting a json string
		certificate := chartverifier.Report{}
		err := yaml.Unmarshal([]byte(outBuf.String()), &certificate)
		require.NoError(t, err)
		require.True(t, certificate.Metadata.ToolMetadata.ProviderDelivery)

	})

	t.Run("should see providerControlledDelivery is false if no -d flag", func(t *testing.T) {

		cmd := NewVerifyCmd(viper.New())
		outBuf := bytes.NewBufferString("")
		utils.CmdStdout = outBuf
		errBuf := bytes.NewBufferString("")
		cmd.SetErr(errBuf)

		cmd.SetArgs([]string{
			"-e", "has-readme", // only consider a single check, perhaps more checks in the future
			"../pkg/chartverifier/checks/chart-0.1.0-v3.valid.tgz",
			"-E"})

		require.NoError(t, cmd.Execute())
		require.NotEmpty(t, outBuf.String())

		// attempts to deserialize the command's output, expecting a json string
		certificate := chartverifier.Report{}
		err := yaml.Unmarshal([]byte(outBuf.String()), &certificate)
		require.NoError(t, err)
		require.False(t, certificate.Metadata.ToolMetadata.ProviderDelivery)

	})

}

func TestBuildChecks(t *testing.T) {
	t.Run("Should fail when enabledChecks and disabledChecks have more than one item at the same time", func(t *testing.T) {
		var (
			all      = make(checks.DefaultRegistry)
			enabled  = []string{string(checks.HasReadmeName)}
			disabled = []string{string(checks.ChartTestingName)}
		)
		all.Add(checks.HasReadmeName, "v1.0", nil)
		all.Add(checks.ChartTestingName, "v1.0", nil)
		all.Add(checks.ContainsTestName, "v1.0", nil)
		selected, err := buildChecks(all, viper.New(), enabled, disabled)
		require.Error(t, err)
		require.Nil(t, selected)
	})

	t.Run("Should fail when enabled check is unknown", func(t *testing.T) {
		var (
			all      = make(checks.DefaultRegistry)
			disabled = []string{}
			enabled  = []string{"d"}
		)
		all.Add(checks.HasReadmeName, "v1.0", nil)
		all.Add(checks.ChartTestingName, "v1.0", nil)
		all.Add(checks.ContainsTestName, "v1.0", nil)
		selected, err := buildChecks(all, viper.New(), enabled, disabled)
		require.Error(t, err)
		require.Nil(t, selected)
	})

	t.Run("Should fail when disabled check is unknown", func(t *testing.T) {
		var (
			all      = make(checks.DefaultRegistry)
			disabled = []string{"e"}
			enabled  = []string{}
		)
		all.Add(checks.HasReadmeName, "v1.0", checks.HasReadme)
		all.Add(checks.ChartTestingName, "v1.0", checks.ChartTesting)
		all.Add(checks.ContainsTestName, "v1.0", checks.ContainsTest)
		selected, err := buildChecks(all, viper.New(), enabled, disabled)
		require.Error(t, err)
		require.Nil(t, selected)
	})

	t.Run("Should return all checks when neither enabled or disabled checks have been informed", func(t *testing.T) {
		var (
			enabled  = []string{}
			disabled = []string{}
			all      = make(checks.DefaultRegistry)
		)
		all.Add(checks.HasReadmeName, "v1.0", nil)
		all.Add(checks.ChartTestingName, "v1.0", nil)
		all.Add(checks.ContainsTestName, "v1.0", nil)
		selected, err := buildChecks(all, viper.New(), enabled, disabled)
		require.NoError(t, err)
		for k, _ := range all {
			_, ok := selected[k.Name]
			require.True(t, ok, fmt.Sprintf("Missing Check: %s", k.Name))
		}
	})

}
