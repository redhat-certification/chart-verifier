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
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/redhat-certification/chart-verifier/internal/chartverifier/junitxml"
	"github.com/redhat-certification/chart-verifier/internal/chartverifier/utils"
	"github.com/redhat-certification/chart-verifier/internal/tool"
	apiChecks "github.com/redhat-certification/chart-verifier/pkg/chartverifier/checks"
	apireport "github.com/redhat-certification/chart-verifier/pkg/chartverifier/report"
	apiverifier "github.com/redhat-certification/chart-verifier/pkg/chartverifier/verifier"
	apiversion "github.com/redhat-certification/chart-verifier/pkg/chartverifier/version"

	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/cli/values"

	//"helm.sh/helm/v3/pkg/getter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	//"github.com/redhat-certification/chart-verifier/internal/chartverifier"
)

//func init() {
//	allChecks = chartverifier.DefaultRegistry().AllChecks()
//}

//goland:noinspection GoUnusedGlobalVariable
var (
	// enabledChecksFlag are the checks that should be performed, after the command initialization has happened.
	enabledChecksFlag []string
	// disabledChecksFlag are the checks that should not be performed.
	disabledChecksFlag []string
	// outputFormatFlag contains the output format the user has specified: default, yaml or json.
	outputFormatFlag string
	// openshiftVersionFlag set the value of `certifiedOpenShiftVersions` in the report
	openshiftVersionFlag string
	// write report to file
	reportToFile bool
	// write an error log file
	suppressErrorLog bool
	// skip helm cleanup
	skipCleanup bool
	// distribution method is web-catalog-only.
	webCatalogOnly bool
	// client timeout
	clientTimeout time.Duration
	// pgp public key file
	pgpPublicKeyFile string
	// helm install timeout
	helmInstallTimeout time.Duration
	// writeJUnitXMLTo is where to write an additional junitxml representation of the outcome
	writeJUnitXMLTo string
)

func buildChecks(enabled []string, unEnabled []string) ([]apiChecks.CheckName, []apiChecks.CheckName, error) {
	var enabledChecks []apiChecks.CheckName
	var unEnabledChecks []apiChecks.CheckName
	var convertErr error
	if len(enabled) > 0 && len(unEnabled) > 0 {
		return enabledChecks, unEnabledChecks, errors.New("--enable and --disable can't be used at the same time")
	} else if len(enabled) > 0 {
		enabledChecks, convertErr = convertChecks(enabled)
		if convertErr != nil {
			return enabledChecks, unEnabledChecks, convertErr
		}
	} else if len(unEnabled) > 0 {
		unEnabledChecks, convertErr = convertChecks(unEnabled)
		if convertErr != nil {
			return enabledChecks, unEnabledChecks, convertErr
		}
	}
	return enabledChecks, unEnabledChecks, nil
}

func convertChecks(checks []string) ([]apiChecks.CheckName, error) {
	var apiCheckSet []apiChecks.CheckName
	for _, check := range checks {
		checkName := apiChecks.CheckName(check)
		checkFound := slices.Contains(apiChecks.GetChecks(), checkName)
		if checkFound {
			apiCheckSet = append(apiCheckSet, checkName)
		} else {
			return apiCheckSet, fmt.Errorf("enabled check is invalid :%s", check)
		}
	}
	return apiCheckSet, nil
}

func convertToMap(values []string) map[string]interface{} {
	valueMap := make(map[string]interface{})
	for _, val := range values {
		parts := strings.SplitN(val, "=", 2)
		valueMap[parts[0]] = parts[1]
	}
	return valueMap
}

// settings comes from Helm, to extract the same configuration values Helm uses.
var settings = cli.New()

type verifyOptions struct {
	ValueFiles []string
	Values     []string
}

// NewVerifyCmd creates ...
func NewVerifyCmd(config *viper.Viper) *cobra.Command {
	// opts contains command line options extracted from the environment.
	opts := &values.Options{}

	// verifyOpts contains this specific command options.
	verifyOpts := &verifyOptions{}

	cmd := &cobra.Command{
		Use:   "verify <chart-uri>",
		Args:  cobra.ExactArgs(1),
		Short: "Verifies a Helm chart by checking some of its characteristics",
		RunE: func(cmd *cobra.Command, args []string) error {
			reportFormat := apireport.YamlReport
			if outputFormatFlag == "json" {
				reportFormat = apireport.JSONReport
			}

			reportName := ""
			if reportToFile {
				if outputFormatFlag == "json" {
					reportName = "report.json"
				} else {
					reportName = "report.yaml"
				}
			}

			enabledChecks, unEnabledChecks, checksErr := buildChecks(enabledChecksFlag, disabledChecksFlag)
			if checksErr != nil {
				return checksErr
			}

			utils.InitLog(cmd, reportName, suppressErrorLog)

			utils.LogInfo(fmt.Sprintf("Chart Verifer %s.", apiversion.GetVersion()))
			utils.LogInfo(fmt.Sprintf("Verify : %s", args[0]))
			utils.LogInfo(fmt.Sprintf("Client timeout: %s", clientTimeout))
			utils.LogInfo(fmt.Sprintf("Helm Install timeout: %s", helmInstallTimeout))

			valueMap := convertToMap(verifyOpts.Values)
			for key, val := range viper.AllSettings() {
				valueMap[strings.ToLower(key)] = val
			}

			verifier := apiverifier.NewVerifier()

			if len(enabledChecks) > 0 {
				verifier = verifier.EnableChecks(enabledChecks)
			} else if len(unEnabledChecks) > 0 {
				verifier = verifier.UnEnableChecks(unEnabledChecks)
			}

			encodedKey, err := tool.GetEncodedKey(pgpPublicKeyFile)
			if err != nil {
				return err
			}

			var runErr error
			verifier, runErr = verifier.SetBoolean(apiverifier.WebCatalogOnly, webCatalogOnly).
				SetBoolean(apiverifier.SuppressErrorLog, suppressErrorLog).
				SetBoolean(apiverifier.SkipCleanup, skipCleanup).
				SetDuration(apiverifier.Timeout, clientTimeout).
				SetDuration(apiverifier.HelmInstallTimeout, helmInstallTimeout).
				SetString(apiverifier.OpenshiftVersion, []string{openshiftVersionFlag}).
				SetString(apiverifier.ChartValues, opts.ValueFiles).
				SetString(apiverifier.KubeAPIServer, []string{settings.KubeAPIServer}).
				SetString(apiverifier.KubeAsUser, []string{settings.KubeAsUser}).
				SetString(apiverifier.KubeCaFile, []string{settings.KubeCaFile}).
				SetString(apiverifier.KubeConfig, []string{settings.KubeConfig}).
				SetString(apiverifier.KubeContext, []string{settings.KubeContext}).
				SetString(apiverifier.Namespace, []string{settings.Namespace()}).
				SetString(apiverifier.KubeAPIServer, []string{settings.KubeAPIServer}).
				SetString(apiverifier.RegistryConfig, []string{settings.RegistryConfig}).
				SetString(apiverifier.RepositoryConfig, []string{settings.RepositoryConfig}).
				SetString(apiverifier.RepositoryCache, []string{settings.RepositoryCache}).
				SetString(apiverifier.KubeAsGroups, settings.KubeAsGroups).
				SetValues(apiverifier.CommandSet, valueMap).
				SetValues(apiverifier.ChartSet, convertToMap(opts.Values)).
				SetValues(apiverifier.ChartSetFile, convertToMap(opts.FileValues)).
				SetValues(apiverifier.ChartSetString, convertToMap(opts.StringValues)).
				SetString(apiverifier.PGPPublicKey, []string{encodedKey}).
				Run(args[0])

			if runErr != nil {
				return runErr
			}

			report, reportErr := verifier.GetReport().GetContent(reportFormat)
			if reportErr != nil {
				return reportErr
			}

			// Failure to write JUnitXML result is non-fatal because junitxml reports are considered extra.
			if writeJUnitXMLTo != "" {
				utils.LogInfo(fmt.Sprintf("user requested additional junitxml report be written to %s", writeJUnitXMLTo))
				junitOutput, err := junitxml.Format(*verifier.GetReport())
				if err != nil {
					utils.LogError(fmt.Sprintf("failed to convert report content to junitxml: %s", err))
				} else {
					err = os.WriteFile(writeJUnitXMLTo, junitOutput, 0o644)
					if err != nil {
						utils.LogError(fmt.Sprintf("failed to write junitxml output to specified path %s: %s", writeJUnitXMLTo, err))
					}
				}
			}

			utils.WriteStdOut(report)

			utils.WriteLogs(outputFormatFlag)

			return nil
		},
	}

	settings.AddFlags(cmd.Flags())

	cmd.Flags().StringSliceVarP(&opts.ValueFiles, "chart-values", "F", nil, "specify values in a YAML file or a URL (can specify multiple)")

	cmd.Flags().StringSliceVarP(&opts.Values, "chart-set", "S", nil, "set values for the chart (can specify multiple or separate values with commas: key1=val1,key2=val2)")

	cmd.Flags().StringSliceVarP(&opts.StringValues, "chart-set-string", "X", nil, "set STRING values for the chart (can specify multiple or separate values with commas: key1=val1,key2=val2)")

	cmd.Flags().StringSliceVarP(&opts.FileValues, "chart-set-file", "G", nil, "set values from respective files specified via the command line (can specify multiple or separate values with commas: key1=path1,key2=path2)")

	cmd.Flags().StringSliceVarP(&enabledChecksFlag, "enable", "e", nil, "only the informed checks will be enabled")

	cmd.Flags().StringSliceVarP(&disabledChecksFlag, "disable", "x", nil, "all checks will be enabled except the informed ones")

	cmd.Flags().StringVarP(&outputFormatFlag, "output", "o", "", "the output format: default, json or yaml")

	cmd.Flags().StringSliceVarP(&verifyOpts.Values, "set", "s", []string{}, "overrides a configuration, e.g: dummy.ok=false")

	cmd.Flags().StringSliceVarP(&verifyOpts.ValueFiles, "set-values", "f", nil, "specify application and check configuration values in a YAML file or a URL (can specify multiple)")
	cmd.Flags().StringVarP(&openshiftVersionFlag, "openshift-version", "V", "", "version of OpenShift used in the cluster")
	cmd.Flags().DurationVar(&clientTimeout, "timeout", 30*time.Minute, "time to wait for completion of chart install and test")
	cmd.Flags().BoolVarP(&reportToFile, "write-to-file", "w", false, "write report to ./chartverifier/report.yaml (default: stdout)")
	cmd.Flags().BoolVarP(&suppressErrorLog, "suppress-error-log", "E", false, "suppress the error log (default: written to ./chartverifier/verifier-<timestamp>.log)")
	cmd.Flags().BoolVarP(&skipCleanup, "skip-cleanup", "c", false, "set this to skip resource cleanup after verifier run")
	cmd.Flags().BoolVarP(&webCatalogOnly, "web-catalog-only", "W", false, "set this to indicate that the distribution method is web catalog only (default: false)")
	cmd.Flags().StringVarP(&pgpPublicKeyFile, "pgp-public-key", "k", "", "file containing gpg public key of the key used to sign the chart")
	cmd.Flags().DurationVar(&helmInstallTimeout, "helm-install-timeout", 5*time.Minute, "helm install timeout")
	cmd.Flags().StringVar(&writeJUnitXMLTo, "write-junitxml-to", "", "If set, will write a junitXML representation of the result to the specified path in addition to the configured output format")

	return cmd
}

func init() {
	rootCmd.AddCommand(NewVerifyCmd(viper.GetViper()))
}
