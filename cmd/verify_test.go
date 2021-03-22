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

	"github.com/redhat-certification/chart-verifier/pkg/chartverifier/checks"
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

		cmd.SetArgs([]string{"../pkg/chartverifier/checks/chart-0.1.0-v3.non-existing.tgz"})

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

		cmd.SetArgs([]string{"../pkg/chartverifier/checks/chart-0.1.0-v3.non-existing.tgz", "-o"})
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
			"../pkg/chartverifier/checks/chart-0.1.0-v3.valid.tgz",
		})
		err := cmd.Execute()
		require.Error(t, err)
		require.False(t, checks.IsChartNotFound(err))
	})

	t.Run("Should succeed when the chart exists and is valid for a single check", func(t *testing.T) {
		cmd := NewVerifyCmd(viper.New())
		outBuf := bytes.NewBufferString("")
		cmd.SetOut(outBuf)
		errBuf := bytes.NewBufferString("")
		cmd.SetErr(errBuf)

		cmd.SetArgs([]string{
			"-e", "is-helm-v3",
			"../pkg/chartverifier/checks/chart-0.1.0-v3.valid.tgz",
		})
		require.NoError(t, cmd.Execute())
		require.NotEmpty(t, outBuf.String())

		expected := "chart: chart\n" +
			"version: 1.16.0\n" +
			"ok: true\n" +
			"\n" +
			"is-helm-v3:\n" +
			"\tok: true\n" +
			"\treason: " + checks.Helm3Reason + "\n"
		require.Equal(t, expected, outBuf.String())
	})

	t.Run("Should display JSON certificate when option --output and argument values are given", func(t *testing.T) {
		cmd := NewVerifyCmd(viper.New())
		outBuf := bytes.NewBufferString("")
		cmd.SetOut(outBuf)
		errBuf := bytes.NewBufferString("")
		cmd.SetErr(errBuf)

		cmd.SetArgs([]string{
			"-e", "is-helm-v3", // only consider a single check, perhaps more checks in the future
			"-o", "json",
			"../pkg/chartverifier/checks/chart-0.1.0-v3.valid.tgz",
		})
		require.NoError(t, cmd.Execute())
		require.NotEmpty(t, outBuf.String())

		// attempts to deserialize the command's output, expecting a json string
		actual := map[string]interface{}{}
		err := json.Unmarshal([]byte(outBuf.String()), &actual)
		require.NoError(t, err)

		expected := map[string]interface{}{
			"metadata": map[string]interface{}{
				"chart": map[string]interface{}{
					"name":    "chart",
					"version": "1.16.0",
				},
			},
			"ok": true,
			"results": map[string]interface{}{
				"is-helm-v3": map[string]interface{}{
					"ok":     true,
					"reason": checks.Helm3Reason,
				},
			},
		}
		require.Equal(t, expected, actual)
	})

	t.Run("Should display YAML certificate when option --output and argument values are given", func(t *testing.T) {
		cmd := NewVerifyCmd(viper.New())
		outBuf := bytes.NewBufferString("")
		cmd.SetOut(outBuf)
		errBuf := bytes.NewBufferString("")
		cmd.SetErr(errBuf)

		cmd.SetArgs([]string{
			"-e", "is-helm-v3", // only consider a single check, perhaps more checks in the future
			"-o", "yaml",
			"../pkg/chartverifier/checks/chart-0.1.0-v3.valid.tgz",
		})
		require.NoError(t, cmd.Execute())
		require.NotEmpty(t, outBuf.String())

		// attempts to deserialize the command's output, expecting a json string
		actual := map[string]interface{}{}
		err := yaml.Unmarshal([]byte(outBuf.String()), &actual)
		require.NoError(t, err)

		expected := map[string]interface{}{
			"metadata": map[string]interface{}{
				"chart": map[string]interface{}{
					"name":    "chart",
					"version": "1.16.0",
				},
			},
			"ok": true,
			"results": map[string]interface{}{
				"is-helm-v3": map[string]interface{}{
					"ok":     true,
					"reason": checks.Helm3Reason,
				},
			},
		}
		require.Equal(t, expected, actual)
	})

}

func TestBuildChecks(t *testing.T) {
	t.Run("Should fail when enabledChecks and disabledChecks have more than one item at the same time", func(t *testing.T) {
		var (
			all      = []string{"a", "b", "c"}
			enabled  = []string{"a"}
			disabled = []string{"b"}
		)
		selected, err := buildChecks(all, enabled, disabled)
		require.Error(t, err)
		require.Nil(t, selected)
	})

	t.Run("Should fail when enabled check is unknown", func(t *testing.T) {
		var (
			all      = []string{"a", "b", "c"}
			disabled = []string{}
			enabled  = []string{"d"}
		)
		selected, err := buildChecks(all, enabled, disabled)
		require.Error(t, err)
		require.Nil(t, selected)
	})

	t.Run("Should fail when disabled check is unknown", func(t *testing.T) {
		var (
			all      = []string{"a", "b", "c"}
			disabled = []string{"e"}
			enabled  = []string{}
		)
		selected, err := buildChecks(all, enabled, disabled)
		require.Error(t, err)
		require.Nil(t, selected)
	})

	t.Run("Should return all checks when neither enabled or disabled checks have been informed", func(t *testing.T) {
		var (
			enabled  = []string{}
			disabled = []string{}
			all      = []string{"a", "b", "c"}
		)
		selected, err := buildChecks(all, enabled, disabled)
		require.NoError(t, err)
		require.Equal(t, selected, all)
	})

}
