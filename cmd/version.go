package cmd

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	apiversion "github.com/redhat-certification/chart-verifier/pkg/chartverifier/version"
)

// Print version and commit ID as json blob
var asData bool

func init() {
	rootCmd.AddCommand(NewVersionCmd())
}

// CommitIDLong contains the commit ID the binary was build on. It is populated at build time by ldflags.
// If you're running from a local debugger it will show an empty commit ID.
var CommitIDLong string = "unknown"

type VersionContext struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
}

var Version = VersionContext{
	Version: apiversion.GetVersion(),
	Commit:  CommitIDLong,
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
	cmd.Flags().BoolVar(&asData, "as-data", false, "output the version and commit ID information in JSON format")
	return cmd
}

func runVersion(out io.Writer) error {
	if asData {
		marshalledVersion, err := json.Marshal(Version)
		if err != nil {
			return err
		}
		fmt.Fprintf(out, "%s\n", string(marshalledVersion))
		return nil
	}

	fmt.Fprintf(out, "chart-verifier v%s <commit: %s>\n", apiversion.GetVersion(), CommitIDLong)
	return nil
}
