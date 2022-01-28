package cmd

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

var Version = "0.0.0"

//go:embed release/release_info.json
var releaseFileContent []byte

type Release struct {
	Version string `json:"version"`
}

func init() {

	var release Release
	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	err := json.Unmarshal(releaseFileContent, &release)
	if err != nil {
		Version = "0.0.3"
		return
	}

	Version = release.Version
	rootCmd.AddCommand(newVersionCmd())
}

func newVersionCmd() *cobra.Command {
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
	if Version == "0.0.0" {
		return fmt.Errorf("no version info available")
	}
	fmt.Fprintf(out, "v%s\n", Version)
	return nil
}
