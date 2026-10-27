package cmd

import (
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"

	"github.com/redhat-certification/chart-verifier/internal/chartverifier/profiles"
)

func init() {
	rootCmd.AddCommand(NewListChecksCmd())
}

func NewListChecksCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list-checks",
		Short: "List all checks that will be executed for each profile",
		Long:  "This command lists all checks that chart-verifier uses against a Helm chart, grouped by vendor profile.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			printChecks(cmd.OutOrStdout())
			return nil
		},
	}
}

func printChecks(w io.Writer) {
	fmt.Fprintln(w, "These are the available checks for each profile:")
	for _, profile := range profiles.ListLatest() {
		fmt.Fprintln(w, formattedProfileBlock(profile))
	}
}

func formattedProfileBlock(profile *profiles.Profile) string {
	title := fmt.Sprintf("[%s Profile %s]: %s", profile.Vendor, profile.Version, profileDescription(profile.Vendor))
	return strings.Join([]string{title, formatCheckList(profile)}, "\n")
}

func profileDescription(vendor profiles.VendorType) string {
	switch vendor {
	case profiles.DefaultProfile:
		return "default for partner Helm chart certification"
	case "redhat":
		return "invoked with --set profile.vendortype=redhat"
	case "community":
		return "invoked with --set profile.vendortype=community"
	case "developer-console":
		return "invoked with --set profile.vendortype=developer-console"
	default:
		return fmt.Sprintf("invoked with --set profile.vendortype=%s", vendor)
	}
}

func formatCheckList(profile *profiles.Profile) string {
	var s string
	for _, check := range profile.Checks {
		s += dashPrefix(formatCheck(check)) + "\n"
	}
	return s
}

func formatCheck(check *profiles.Check) string {
	name := check.Name
	version := ""
	if parts := strings.SplitN(check.Name, "/", 2); len(parts) == 2 {
		version = parts[0]
		name = parts[1]
	}
	if version != "" {
		return fmt.Sprintf("%s (%s, %s)", name, version, check.Type)
	}
	return fmt.Sprintf("%s (%s)", name, check.Type)
}

func dashPrefix(s string) string {
	return fmt.Sprintf("- %s", s)
}
