package providers

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Run("It should creates default values", func(t *testing.T) {
		c, err := NewConfig()
		require.NoError(t, err)
		require.Equal(t, &Config{
			DatabaseUser:     defaultDbUser,
			DatabasePassword: defaultDbPassword,
			DatabaseName:     defaultDbName,
			ListenAddr:       defaultAddr,
			HashComplexity:   defaultHashComplexity,
		}, c)
	})

	t.Run("All values should be overridable", func(t *testing.T) {
		os.Setenv("DB_USER", "DB_USER")
		os.Setenv("DB_PASSWORD", "DB_PASSWORD")
		os.Setenv("DB_NAME", "DB_NAME")
		os.Setenv("ADDR", "ADDR")

		c, err := NewConfig()
		require.NoError(t, err)
		require.Equal(t, &Config{
			DatabaseUser:     "DB_USER",
			DatabasePassword: "DB_PASSWORD",
			DatabaseName:     "DB_NAME",
			ListenAddr:       "ADDR",
			HashComplexity:   defaultHashComplexity,
		}, c)
	})
}
