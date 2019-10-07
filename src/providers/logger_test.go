package providers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	logger, err := NewLogger()
	require.NoError(t, err)
	require.NotNil(t, logger)
}
