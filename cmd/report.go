package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/redhat-certification/chart-verifier/internal/chartverifier/utils"
	apireport "github.com/redhat-certification/chart-verifier/pkg/chartverifier/report"
	apireportsummary "github.com/redhat-certification/chart-verifier/pkg/chartverifier/reportsummary"
)

func init() {
	rootCmd.AddCommand(NewReportCmd(viper.GetViper()))
}

type reportOptions struct {
	ValueFiles []string
	Values     []string
}

var skipDigestCheck bool

// NewReportCmd creates a command that sanity checks report.
func NewReportCmd(config *viper.Viper) *cobra.Command {
	// verifyOpts contains this specific command options.
	reportOpts := &reportOptions{}

	cmd := &cobra.Command{
		Use: fmt.Sprintf("report {%s,%s,%s,%s,%s} <report-uri>", apireportsummary.AllSummary, apireportsummary.AnnotationsSummary, apireportsummary.DigestsSummary,
			apireportsummary.MetadataSummary, apireportsummary.ResultsSummary),
		Args:  cobra.ExactArgs(2),
		Short: "Provides information from a report",
		RunE: func(cmd *cobra.Command, args []string) error {
			reportName := ""
			reportFormat := apireportsummary.JSONReport
			if reportToFile {
				if outputFormatFlag == "json" {
					reportName = "report-info.json"
				} else {
					reportName = "report-info.yaml"
					reportFormat = apireportsummary.YAMLReport
				}
			}
			utils.InitLog(cmd, reportName, true)

			commandArg := args[0]
			reportArg := args[1]

			var reportType apireportsummary.SummaryType
			switch commandArg {
			case string(apireportsummary.MetadataSummary):
				reportType = apireportsummary.MetadataSummary
			case string(apireportsummary.DigestsSummary):
				reportType = apireportsummary.DigestsSummary
			case string(apireportsummary.AnnotationsSummary):
				reportType = apireportsummary.AnnotationsSummary
			case string(apireportsummary.ResultsSummary):
				reportType = apireportsummary.ResultsSummary
			case string(apireportsummary.AllSummary):
				reportType = apireportsummary.AllSummary
			default:
				return fmt.Errorf("Error: command %s not recognized", commandArg)
			}

			valueMap := make(map[string]interface{})
			for _, val := range reportOpts.Values {
				parts := strings.Split(val, "=")
				valueMap[parts[0]] = parts[1]
			}
			for key, val := range viper.AllSettings() {
				valueMap[key] = val
			}
			// #nosec G304
			reportFile, openErr := os.Open(reportArg)
			if openErr != nil {
				return errors.New(fmt.Sprintf("report path %s: error opening file  %v", reportArg, openErr))
			}

			reportBytes, readErr := io.ReadAll(reportFile)
			if readErr != nil {
				//nolint:errcheck // TODO(komish): The linter indicates that we've not done anything with
				// this error, which is correct. Need to confirm what the intention was before changing this.
				fmt.Errorf("report path %s: error reading file  %v", reportArg, readErr)
			}

			report, loadErr := apireport.NewReport().
				SetContent(string(reportBytes)).
				Load()
			if loadErr != nil {
				return loadErr
			}

			reportSummary, summaryErr := apireportsummary.NewReportSummary().
				SetValues(valueMap).
				SetReport(report).
				SetBoolean(apireportsummary.SkipDigestCheck, skipDigestCheck).
				GetContent(reportType, reportFormat)

			if summaryErr != nil {
				return errors.New(fmt.Sprintf("Error executing command: %v", summaryErr))
			}

			utils.WriteStdOut(reportSummary)
			return nil
		},
	}

	cmd.Flags().StringVarP(&outputFormatFlag, "output", "o", "", "the output format: json (default) or yaml")

	cmd.Flags().StringSliceVarP(&reportOpts.Values, "set", "s", []string{}, "set report configuration values: profile vendor type and version")

	cmd.Flags().StringSliceVarP(&reportOpts.ValueFiles, "set-values", "f", nil, "specify report configuration values in a YAML file or a URL (can specify multiple)")

	cmd.Flags().BoolVarP(&reportToFile, "write-to-file", "w", false, "write report to report-info.json (default: stdout)")

	cmd.Flags().BoolVarP(&skipDigestCheck, "skip-digest-check", "d", false, "FOR TESTING PURPOSES ONLY: skip the check that the digest in the report matches the report content")

	return cmd
}
