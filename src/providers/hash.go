package providers

import (
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"golang.org/x/crypto/bcrypt"
)

type (
	// Hasher is interface for implementing hashing and validating password
	Hasher interface {
		Hash(password string) (string, error)
		Compare(password, hash string) bool
	}

	// BcryptHasher default bcrypt that implements
	BcryptHasher struct {
		complexity int
	}

	// BcryptHasherOptions option for bcrypt hashing algos
	BcryptHasherOptions struct {
		fx.In

		Config     *viper.Viper
		Complexity int
	}
)

// NewBcryptHasher creates a new hasher that uses bcrypt under the hood
func NewBcryptHasher(config *Config) Hasher {
	return &BcryptHasher{
		complexity: config.HashComplexity,
	}
}

// Hash hash given password with complexety using bcrypt
func (b BcryptHasher) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), b.complexity)
	return string(bytes), err
}

// Compare compare password to hash
func (b BcryptHasher) Compare(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
