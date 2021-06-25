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
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestAllPassReport(t *testing.T) {
	t.Run("Default profile with all pass result", func(t *testing.T) {
		cmd := NewCheckCmd(viper.New())
		outBuf := bytes.NewBufferString("")
		cmd.SetOut(outBuf)
		errBuf := bytes.NewBufferString("")
		cmd.SetErr(errBuf)

		cmd.SetArgs([]string{
			"testdata/report-all-pass.yaml",
		})
		require.NoError(t, cmd.Execute())
		require.NotEmpty(t, outBuf.String())

		expected := "profile: default\n" +
			"passed: 11\n" +
			"failed: 0\n" +
			"unknown: 0\n" +
			"message: []\n"
		require.Contains(t, outBuf.String(), expected)
	})
}

func TestInvalidReport(t *testing.T) {
	t.Run("Default profile with invalid result", func(t *testing.T) {
		cmd := NewCheckCmd(viper.New())
		outBuf := bytes.NewBufferString("")
		cmd.SetOut(outBuf)
		errBuf := bytes.NewBufferString("")
		cmd.SetErr(errBuf)

		cmd.SetArgs([]string{
			"testdata/report-invalid.yaml",
		})
		require.Error(t, cmd.Execute())
		require.NotEmpty(t, outBuf.String())
		require.NotEmpty(t, errBuf.String())

		expectedStdOut :=
			`Usage:
  check <report-uri> [flags]

Flags:
      --debug                       enable verbose output
  -h, --help                        help for check
      --kube-apiserver string       the address and the port for the Kubernetes API server
      --kube-as-group stringArray   group to impersonate for the operation, this flag can be repeated to specify multiple groups.
      --kube-as-user string         username to impersonate for the operation
      --kube-ca-file string         the certificate authority file for the Kubernetes API server connection
      --kube-context string         name of the kubeconfig context to use
      --kube-token string           bearer token used for authentication
      --kubeconfig string           path to the kubeconfig file
  -n, --namespace string            namespace scope for this request
  -o, --output string               the output format: default, json or yaml
  -p, --profile string              check according to specific profile (partner, redhat, community) (default "default")
      --registry-config string      path to the registry config file (default "/home/abai/.config/helm/registry.json")
      --repository-cache string     path to the file containing cached repository indexes (default "/home/abai/.cache/helm/repository")
      --repository-config string    path to the file containing repository names and URLs (default "/home/abai/.config/helm/repositories.yaml")
`
		expectedStdErr := `checking report results: incorrect outcome type 'INVALID' for 'chart-testing'`
		require.Contains(t, outBuf.String(), expectedStdOut)
		require.Contains(t, errBuf.String(), expectedStdErr)
	})
}

func TestReportWithFails(t *testing.T) {
	t.Run("Default profile with mixed outcome type", func(t *testing.T) {
		cmd := NewCheckCmd(viper.New())
		outBuf := bytes.NewBufferString("")
		cmd.SetOut(outBuf)
		errBuf := bytes.NewBufferString("")
		cmd.SetErr(errBuf)

		cmd.SetArgs([]string{
			"testdata/report-with-fails.yaml",
		})
		require.NoError(t, cmd.Execute())
		require.NotEmpty(t, outBuf.String())

		expected :=
			`profile: default
passed: 8
failed: 2
unknown: 1
message:
  - check: images-are-certified
    type: Mandatory
    outcome: FAIL
    reason: 'Error: some images are not certified.'
  - check: helm-lint
    type: Mandatory
    outcome: UNKNOWN
    reason: Unknown outcome.
  - check: chart-testing
    type: Mandatory
    outcome: FAIL
    reason: 'Error: chart testin failed.'

`
		require.Contains(t, outBuf.String(), expected)
	})
}
