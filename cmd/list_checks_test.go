package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/redhat-certification/chart-verifier/internal/chartverifier/profiles"
	apiChecks "github.com/redhat-certification/chart-verifier/pkg/chartverifier/checks"
)

func TestFormatCheck(t *testing.T) {
	require.Equal(t, "helm-lint (v1.0, Mandatory)", formatCheck(&profiles.Check{
		Name: "v1.0/helm-lint",
		Type: apiChecks.MandatoryCheckType,
	}))
	require.Equal(t, "has-readme (Optional)", formatCheck(&profiles.Check{
		Name: "has-readme",
		Type: apiChecks.OptionalCheckType,
	}))
}

func TestDashPrefix(t *testing.T) {
	require.True(t, strings.HasPrefix(dashPrefix("foo"), "- "))
	require.Equal(t, "- foo", dashPrefix("foo"))
}

func TestFormatCheckList(t *testing.T) {
	profile := &profiles.Profile{
		Checks: []*profiles.Check{
			{Name: "v1.0/helm-lint", Type: apiChecks.MandatoryCheckType},
			{Name: "v1.0/has-notes", Type: apiChecks.OptionalCheckType},
		},
	}
	got := formatCheckList(profile)
	lines := strings.Split(got, "\n")
	require.Equal(t, len(profile.Checks)+1, len(lines))
	require.Contains(t, got, "- helm-lint (v1.0, Mandatory)")
	require.Contains(t, got, "- has-notes (v1.0, Optional)")
}

func TestFormattedProfileBlock(t *testing.T) {
	profile := &profiles.Profile{
		Vendor:  profiles.DefaultProfile,
		Version: "v1.3",
		Checks: []*profiles.Check{
			{Name: "v1.0/helm-lint", Type: apiChecks.MandatoryCheckType},
		},
	}
	got := formattedProfileBlock(profile)
	require.Contains(t, got, "[partner Profile v1.3]: default for partner Helm chart certification")
	require.Contains(t, got, "- helm-lint (v1.0, Mandatory)")
}

func TestPrintChecks(t *testing.T) {
	buf := &bytes.Buffer{}
	printChecks(buf)
	out := buf.String()

	require.Contains(t, out, "These are the available checks for each profile:")
	require.Contains(t, out, "[partner Profile")
	require.Contains(t, out, "[redhat Profile")
	require.Contains(t, out, "[community Profile")
	require.Contains(t, out, "[developer-console Profile")
	require.Contains(t, out, "helm-lint")
	require.Contains(t, out, "chart-testing")
	require.Contains(t, out, "Mandatory")
	require.Contains(t, out, "Optional")
}

func TestListChecksCommand(t *testing.T) {
	cmd := NewListChecksCmd()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	require.NoError(t, err)

	expected := &bytes.Buffer{}
	printChecks(expected)
	require.Equal(t, expected.String(), buf.String())
}
