package samples

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddStringReport(t *testing.T) {
	err := runVerifier()
	require.NoError(t, err)
}
