package cmd

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"golang.org/x/mod/semver"
	"testing"
)

func TestVersion(t *testing.T) {

	t.Run("Check Version is set.", func(t *testing.T) {
		fmt.Println(fmt.Sprintf("Version is %s", Version))
		require.True(t, semver.IsValid("v"+Version), fmt.Sprintf("Version is not a valid semantic version: %s", Version))
		require.True(t, semver.Compare("v"+Version, "v.0.0.3") > 0, fmt.Sprintf("Version has not been set: %s", Version))
	})

}
