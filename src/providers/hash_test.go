package providers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBcryptHashHash(t *testing.T) {
	hasher := NewBcryptHasher(&Config{
		HashComplexity: 2,
	})

	tests := []string{"123", "password", "new!"}
	for _, test := range tests {
		hash, err := hasher.Hash(test)
		require.NoError(t, err)
		require.True(t, hasher.Compare(test, hash))
	}
}
