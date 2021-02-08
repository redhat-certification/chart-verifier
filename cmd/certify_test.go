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

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/redhat-certification/chart-verifier/pkg/chartverifier/checks"
)

func TestCertify(t *testing.T) {

	t.Run("uri flag is required", func(t *testing.T) {

		t.Run("Should fail when flag -u not given", func(t *testing.T) {
			cmd := NewCertifyCmd()
			outBuf := bytes.NewBufferString("")
			cmd.SetOut(outBuf)
			errBuf := bytes.NewBufferString("")
			cmd.SetErr(errBuf)

			require.Error(t, cmd.Execute())
		})

		t.Run("Should fail when flag -u is given but no value is informed", func(t *testing.T) {
			cmd := NewCertifyCmd()
			outBuf := bytes.NewBufferString("")
			cmd.SetOut(outBuf)
			errBuf := bytes.NewBufferString("")
			cmd.SetErr(errBuf)

			cmd.SetArgs([]string{"-u"})
			require.Error(t, cmd.Execute())
		})

		t.Run("Should fail when flag -u and values are given but resource does not exist", func(t *testing.T) {
			cmd := NewCertifyCmd()
			outBuf := bytes.NewBufferString("")
			cmd.SetOut(outBuf)
			errBuf := bytes.NewBufferString("")
			cmd.SetErr(errBuf)

			cmd.SetArgs([]string{"-u", "../pkg/chartverifier/checks/chart-0.1.0-v3.non-existing.tgz"})

			err := cmd.Execute()
			require.Error(t, err)
			require.True(t, checks.IsChartNotFound(err))
		})

		t.Run("Should fail when flag -o is given but check doesn't exist", func(t *testing.T) {
			cmd := NewCertifyCmd()
			outBuf := bytes.NewBufferString("")
			cmd.SetOut(outBuf)
			errBuf := bytes.NewBufferString("")
			cmd.SetErr(errBuf)

			cmd.SetArgs([]string{"-u", "/tmp/chart.tgz", "-o"})
			err := cmd.Execute()
			require.Error(t, err)
			require.False(t, checks.IsChartNotFound(err))
		})

		t.Run("Should succeed when flag -u and values are given", func(t *testing.T) {
			cmd := NewCertifyCmd()
			outBuf := bytes.NewBufferString("")
			cmd.SetOut(outBuf)
			errBuf := bytes.NewBufferString("")
			cmd.SetErr(errBuf)

			cmd.SetArgs([]string{
				"-u", "../pkg/chartverifier/checks/chart-0.1.0-v3.valid.tgz",
				"--only", "is-helm-v3", // only consider a single check, perhaps more checks in the future
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

		t.Run("Should display JSON certificate when flag --output and -u and values are given", func(t *testing.T) {
			cmd := NewCertifyCmd()
			outBuf := bytes.NewBufferString("")
			cmd.SetOut(outBuf)
			errBuf := bytes.NewBufferString("")
			cmd.SetErr(errBuf)

			cmd.SetArgs([]string{
				"-u", "../pkg/chartverifier/checks/chart-0.1.0-v3.valid.tgz",
				"--only", "is-helm-v3", // only consider a single check, perhaps more checks in the future
				"--output", "json",
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

		t.Run("Should display YAML certificate when flag --output and -u and values are given", func(t *testing.T) {
			cmd := NewCertifyCmd()
			outBuf := bytes.NewBufferString("")
			cmd.SetOut(outBuf)
			errBuf := bytes.NewBufferString("")
			cmd.SetErr(errBuf)

			cmd.SetArgs([]string{
				"-u", "../pkg/chartverifier/checks/chart-0.1.0-v3.valid.tgz",
				"--only", "is-helm-v3", // only consider a single check, perhaps more checks in the future
				"--output", "yaml",
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
	})
}
