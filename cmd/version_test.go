package cmd

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/mod/semver"
)

func TestVersion(t *testing.T) {

	t.Run("Check Version is set.", func(t *testing.T) {
		fmt.Printf("Version is %s", Version)
		require.True(t, semver.IsValid("v"+Version), fmt.Sprintf("Version is not a valid semantic version: %s", Version))
		require.True(t, semver.Compare("v"+Version, "v.0.0.3") > 0, fmt.Sprintf("Version has not been set: %s", Version))
	})

}

func TestVersionCmd(t *testing.T) {
	tests := []struct {
		version   string
		expected  string
		wantError bool
	}{
		{
			version:   "0.0.0",
			expected:  "no version info available",
			wantError: true,
		},
		{
			version:   "1.0.0",
			expected:  "v1.0.0\n",
			wantError: false,
		},
	}
	for _, tt := range tests {
		buf := new(bytes.Buffer)
		Version = tt.version
		err := runVersion(buf)
		if tt.wantError {
			if err == nil {
				t.Errorf("Expected error %q, got none", tt.expected)
			}
			if err.Error() != tt.expected {
				t.Errorf("Expected error %q, got %q", tt.expected, err.Error())
			}
		} else if err != nil {
			t.Errorf("Unexpected error: %s", err)
		} else {
			if buf.String() != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, buf.String())
			}
		}
	}
}
