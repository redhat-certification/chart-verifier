package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

var Version = "0.0.0"

type Release struct {
	Version string `json:"version"`
}

func init() {

	var configDir string
	if isRunningInDockerContainer() {
		configDir = filepath.Join("/app", "releases")
	} else {
		_, fn, _, ok := runtime.Caller(0)
		if !ok {
			return
		}
		index := strings.LastIndex(fn, "chart-verifier/")
		configDir = fn[0 : index+len("chart-verifier")]
		configDir = filepath.Join(configDir, "cmd", "release")
	}

	jsonFile, err := os.Open(filepath.Join(configDir, "release_info.json"))
	if err != nil {
		Version = "0.0.1"
		return
	}

	// read our opened jsonFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)
	if err != nil {
		Version = "0.0.2"
		return
	}

	var release Release
	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	err = json.Unmarshal(byteValue, &release)
	if err != nil {
		Version = "0.0.3"
		return
	}

	Version = release.Version
	rootCmd.AddCommand(newVersionCmd())
}

func isRunningInDockerContainer() bool {
	// docker creates a .dockerenv file at the root
	// of the directory tree inside the container.
	// if this file exists then verifier is running
	// from inside a container
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}
	return false
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
