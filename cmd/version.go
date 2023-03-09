package cmd

import (
	_ "embed"
	"fmt"
	"io"
	"os"

	apiversion "github.com/redhat-certification/chart-verifier/pkg/chartverifier/version"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(NewVersionCmd())
}

func NewVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "print the chart-verifier version information",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runVersion(os.Stdout)
		},
	}
	return cmd
}

func runVersion(out io.Writer) error {
	fmt.Fprintf(out, "v%s\n", apiversion.GetVersion())
	return nil
}
