package cmd

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/redhat-certification/chart-verifier/internal/chartverifier/profiles"
)

// Vendor types that are not useful to print: the default profile is an alias of
// DefaultProfile, and developer-console is not used to certify charts.
var hiddenVendorTypes = map[profiles.VendorType]struct{}{
	profiles.VendorTypeDefault:      {},
	profiles.VendorTypeNotSpecified: {},
	"developer-console":             {},
}

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
	latest := profiles.ListLatest()
	fmt.Fprintln(w, "These are the available checks for each profile:")
	for _, vendor := range listedVendors(latest) {
		fmt.Fprintln(w, formattedProfileBlock(latest[vendor]))
	}
}

// listedVendors returns the vendor types to print, the default profile first and
// the remainder in alphabetical order.
func listedVendors(latest map[profiles.VendorType]*profiles.Profile) []profiles.VendorType {
	vendors := make([]profiles.VendorType, 0, len(latest))
	for vendor := range latest {
		if _, hidden := hiddenVendorTypes[vendor]; !hidden {
			vendors = append(vendors, vendor)
		}
	}
	sort.Slice(vendors, func(i, j int) bool {
		if (vendors[i] == profiles.DefaultProfile) != (vendors[j] == profiles.DefaultProfile) {
			return vendors[i] == profiles.DefaultProfile
		}
		return vendors[i] < vendors[j]
	})
	return vendors
}

func formattedProfileBlock(profile *profiles.Profile) string {
	description := fmt.Sprintf("invoked with --set profile.vendortype=%s", profile.Vendor)
	if profile.Vendor == profiles.DefaultProfile {
		description = "default for Helm chart certification"
	}
	title := fmt.Sprintf("[%s Profile %s]: %s", profile.Vendor, profile.Version, description)
	return strings.Join([]string{title, formatCheckList(profile)}, "\n")
}

func formatCheckList(profile *profiles.Profile) string {
	var s string
	for _, check := range profile.Checks {
		s += fmt.Sprintf("- %s\n", formatCheck(check))
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
