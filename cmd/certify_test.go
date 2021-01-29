/*
 * Copyright (C) 28/12/2020, 16:13 igors
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

package cmd

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"helmcertifier/pkg/helmcertifier/checks"
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

			cmd.SetArgs([]string{"-u", "../pkg/helmcertifier/checks/chart-0.1.0-v3.non-existing.tgz"})

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
				"-u", "../pkg/helmcertifier/checks/chart-0.1.0-v3.valid.tgz",
				"--only", "is-helm-v3", // only consider a single check, perhaps more checks in the future
			})
			require.NoError(t, cmd.Execute())
			require.NotEmpty(t, outBuf.String())

			// FIXME: the chart name inside the tarball should correspond to the tarball name
			expected := "chart: testchart\n" +
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
				"-u", "../pkg/helmcertifier/checks/chart-0.1.0-v3.valid.tgz",
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
						// FIXME: the chart name inside the tarball should correspond to the tarball name
						//"name":    "chart",
						//"version": "0.1.0-v3.valid",
						"name":    "testchart", // should be "chart"
						"version": "1.16.0",    // should be "0.1.0-v3.valid"
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
				"-u", "../pkg/helmcertifier/checks/chart-0.1.0-v3.valid.tgz",
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
						// FIXME: the chart name inside the tarball should correspond to the tarball name
						//"name":    "chart",
						//"version": "0.1.0-v3.valid",
						"name":    "testchart", // should be "chart"
						"version": "1.16.0",    // should be "0.1.0-v3.valid"
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
