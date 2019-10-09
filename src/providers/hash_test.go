package providers

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestBcryptHashHash(t *testing.T) {
	v := viper.New()
	v.Set("hash.complexity", 2)

	hasher := NewBcryptHasher(v)

	tests := []string{"123", "password", "new!"}
	for _, test := range tests {
		hash, err := hasher.Hash(test)
		require.NoError(t, err)
		require.True(t, hasher.Compare(test, hash))
	}
}
