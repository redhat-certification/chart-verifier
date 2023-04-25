package cmd

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/mod/semver"

	apiversion "github.com/redhat-certification/chart-verifier/pkg/chartverifier/version"
)

func TestVersion(t *testing.T) {
	t.Run("Check Version is set.", func(t *testing.T) {
		fmt.Printf("Version is %s", apiversion.GetVersion())
		require.True(t, semver.IsValid("v"+apiversion.GetVersion()), fmt.Sprintf("Version is not a valid semantic version: %s", apiversion.GetVersion()))
		require.True(t, semver.Compare("v"+apiversion.GetVersion(), "v.0.0.3") > 0, fmt.Sprintf("Version has not been set: %s", apiversion.GetVersion()))
	})
}
