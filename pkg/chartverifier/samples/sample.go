package samples

import (
	"fmt"

	"github.com/redhat-certification/chart-verifier/pkg/chartverifier/checks"
	"github.com/redhat-certification/chart-verifier/pkg/chartverifier/report"
	"github.com/redhat-certification/chart-verifier/pkg/chartverifier/reportsummary"
	"github.com/redhat-certification/chart-verifier/pkg/chartverifier/verifier"
)

func runVerifier() error {
	// Run verify command for a chart, but omit the chart testing check and run checks based on the redhat profile
	commandSet := make(map[string]interface{})
	commandSet["profile.vendortype"] = "redhat"

	verifier, verifierErr := verifier.NewVerifier().
		SetValues(verifier.CommandSet, commandSet).
		UnEnableChecks([]checks.CheckName{checks.ChartTesting}).
		Run("https://github.com/redhat-certification/chart-verifier/blob/main/tests/charts/psql-service/0.1.9/psql-service-0.1.9.tgz?raw=true")

	if verifierErr != nil {
		return verifierErr
	}

	// Get and print the report from the verify command
	report, reportErr := verifier.GetReport().
		GetContent(report.YamlReport)
	if reportErr != nil {
		return reportErr
	}
	fmt.Println("report content:\n", report)

	// Get and print the report summary  of the report, but using the partnet profile.
	values := make(map[string]interface{})
	values["profile.vendortype"] = "redhat"

	reportSummary, summmaryErr := reportsummary.NewReportSummary().
		SetReport(verifier.GetReport()).
		SetValues(values).
		GetContent(reportsummary.AllSummary, reportsummary.JsonReport)

	if summmaryErr != nil {
		return summmaryErr
	}
	fmt.Println("report summary content:\n", reportSummary)

	return nil
}
