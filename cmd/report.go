package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/redhat-certification/chart-verifier/pkg/chartverifier/report"
	"github.com/redhat-certification/chart-verifier/pkg/chartverifier/utils"
)

func init() {
	rootCmd.AddCommand(NewReportCmd(viper.GetViper()))
	allSubCommands = report.ReportCommandRegistry().AllCommands()
}

var allSubCommands report.DefaultReportRegistry

type reportOptions struct {
	ValueFiles []string
	Values     []string
}

// NewReportCmd creates a command that sanity checks report.
func NewReportCmd(config *viper.Viper) *cobra.Command {

	// verifyOpts contains this specific command options.
	reportOpts := &reportOptions{}

	cmd := &cobra.Command{
		Use: fmt.Sprintf("report {%s,%s,%s,%s,%s} <report-uri>", report.AllCommandsName, report.AnnotationsCommandName, report.DigestsCommandName,
			report.MetadataCommandName, report.ResultsCommandName),
		Args:  cobra.ExactArgs(2),
		Short: "Provides information from a report",
		RunE: func(cmd *cobra.Command, args []string) error {

			reportName := ""
			if reportToFile {
				if outputFormatFlag == "json" {
					reportName = "report-info.json"
				} else {
					reportName = "report-info.yaml"
				}
			}
			utils.InitLog(cmd, reportName, true)

			commandArg := args[0]
			reportArg := args[1]

			subCommand := report.ReportCommandRegistry().Get(commandArg)

			if subCommand == nil {
				return errors.New(fmt.Sprintf("Error: command %s not recognized", commandArg))
			}

			commandOptions := &report.ReportOptions{}
			commandOptions.AddURI(reportArg)
			commandOptions.AddConfig(config)
			commandOptions.AddValues(reportOpts.Values)

			result, err := subCommand(commandOptions)

			if err != nil {
				return errors.New(fmt.Sprintf("Error executing command: %v", err))
			}

			output := ""
			if outputFormatFlag == "yaml" {
				b, err := yaml.Marshal(result)
				if err != nil {
					utils.LogError(err.Error())
					return err
				}
				output = string(b)
			} else {
				b, err := json.Marshal(result)
				if err != nil {
					utils.LogError(err.Error())
					return err
				}
				output = string(b)
			}
			utils.WriteStdOut(output)

			return nil
		},
	}

	cmd.Flags().StringVarP(&outputFormatFlag, "output", "o", "", "the output format: json (default) or yaml")

	cmd.Flags().StringSliceVarP(&reportOpts.Values, "set", "s", []string{}, "set report configuration values: profile vendor type and version")

	cmd.Flags().StringSliceVarP(&reportOpts.ValueFiles, "set-values", "f", nil, "specify report configuration values in a YAML file or a URL (can specify multiple)")

	cmd.Flags().BoolVarP(&reportToFile, "write-to-file", "w", false, "write report to report-info.json (default: stdout)")

	return cmd
}
