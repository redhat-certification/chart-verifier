package cmd

import (
	"bytes"
	"github.com/redhat-certification/chart-verifier/pkg/chartverifier/report"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestReport(t *testing.T) {

	t.Run("Should fail when no argument is given", func(t *testing.T) {
		cmd := NewReportCmd(viper.New())
		outBuf := bytes.NewBufferString("")
		cmd.SetOut(outBuf)
		errBuf := bytes.NewBufferString("")
		cmd.SetErr(errBuf)

		require.Error(t, cmd.Execute())
	})

	t.Run("Should fail when one argument is given", func(t *testing.T) {
		cmd := NewReportCmd(viper.New())
		outBuf := bytes.NewBufferString("")
		cmd.SetOut(outBuf)
		errBuf := bytes.NewBufferString("")
		cmd.SetErr(errBuf)


		cmd.SetArgs([]string{
			"test/report.yaml",
		})

		require.Error(t, cmd.Execute())
	})

	t.Run("Should fail when bad subcommand is given", func(t *testing.T) {
		cmd := NewReportCmd(viper.New())
		outBuf := bytes.NewBufferString("")
		cmd.SetOut(outBuf)
		errBuf := bytes.NewBufferString("")
		cmd.SetErr(errBuf)

		cmd.SetArgs([]string{
			"None",
			"test/report.yaml",
		})
		require.Error(t, cmd.Execute())
	})

	t.Run("Should pass for subcommand annotations", func(t *testing.T) {
		cmd := NewReportCmd(viper.New())
		outBuf := bytes.NewBufferString("")
		cmd.SetOut(outBuf)
		errBuf := bytes.NewBufferString("")
		cmd.SetErr(errBuf)

		cmd.SetArgs([]string{
			"annotations",
			"test/report.yaml",
		})
		result,err :=  cmd.Execute()
		require.True(t, err==nil, "report annotations failure : %v",err )

		expectedReport := report.OutputReport{}Nm
		expectedReport.AnnotationsReport = append(expectedReport.AnnotationsReport,expectedReport.AnnotationsReport{}

	})
}
