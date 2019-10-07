package services

import (
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/playneta/go-sessions/src/models"
	mock_repositories "github.com/playneta/go-sessions/src/repositories/mocks"
	mock_providers "github.com/playneta/go-sessions/src/providers/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestAccountService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Creating new account repository mock
	account := mock_repositories.NewMockUser(ctrl)

	// Creating new hasher mock
	hasher := mock_providers.NewMockHasher(ctrl)

	// Noop logger
	logger := zap.NewNop().Sugar()

	// Service
	accountService := NewAccount(AccountOptions{
		AccountRepo: account,
		Logger:      logger,
		Hasher:      hasher,
	})

	t.Run("Register", func(t *testing.T) {
		t.Run("Errors", func(t *testing.T) {
			t.Run("Malformed email", func(t *testing.T) {
				user, err := accountService.Register("i am not email", "123456")
				require.Nil(t, user)
				require.Error(t, err, ErrMalformedEmail)
			})

			t.Run("Password not strong enough", func(t *testing.T) {
				user, err := accountService.Register("user@example.com", "123")
				require.Nil(t, user)
				require.Error(t, err, ErrPasswordToSmall)
			})

			t.Run("Password hash failure", func(t *testing.T) {
				hasher.EXPECT().Hash("123456").Return("", errors.New("unknown error"))

				user, err := accountService.Register("user@example.com", "123456")
				require.Nil(t, user)
				require.Error(t, err, errors.New("unknown error"))
			})
		})

		t.Run("Success", func(t *testing.T) {
			newUser := &models.User{
				ID:        1,
				Email:     "user@example.com",
				Password:  "my_password_hash",
				Token:     "",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			account.EXPECT().Create("user@example.com", "my_password_hash").Return(newUser, nil)
			hasher.EXPECT().Hash("123456").Return("my_password_hash", nil)

			user, err := accountService.Register("user@example.com", "123456")
			require.NoError(t, err)
			require.NotNil(t, user)
			require.Equal(t, newUser, user)
		})
	})

	t.Run("Authorize", func(t *testing.T) {
		user := &models.User{
			ID:        1,
			Email:     "user@example.com",
			Password:  "my_password_hash",
			Token:     "",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		t.Run("Errors", func(t *testing.T) {
			t.Run("User not found should return error", func(t *testing.T) {
				account.EXPECT().FindByEmail("user@example.com").Return(nil, nil)

				user, err := accountService.Authorize("user@example.com", "123456")
				require.Nil(t, user)
				require.Error(t, err, ErrUnauthorized)
			})

			t.Run("Database error should return error", func(t *testing.T) {
				account.EXPECT().FindByEmail("user@example.com").Return(nil, errors.New("database error!"))

				user, err := accountService.Authorize("user@example.com", "123456")
				require.Nil(t, user)
				require.Error(t, err, ErrUnauthorized)
			})

			t.Run("Bad password should return error", func(t *testing.T) {
				account.EXPECT().FindByEmail("user@example.com").Return(user, nil)
				hasher.EXPECT().Compare("123456", "my_password_hash").Return(false)

				user, err := accountService.Authorize("user@example.com", "123456")
				require.Nil(t, user)
				require.Error(t, err, ErrUnauthorized)
			})

			t.Run("Update error should return error", func(t *testing.T) {
				account.EXPECT().FindByEmail("user@example.com").Return(user, nil)
				hasher.EXPECT().Compare("123456", "my_password_hash").Return(true)
				account.EXPECT().UpdateToken(user, "secure_token").Return(errors.New("token error!"))

				user, err := accountService.Authorize("user@example.com", "123456")
				require.Nil(t, user)
				require.Error(t, err, errors.New("token error!"))
			})
		})

		t.Run("Success", func(t *testing.T) {
			account.EXPECT().FindByEmail("user@example.com").Return(user, nil)
			hasher.EXPECT().Compare("123456", "my_password_hash").Return(true)
			account.EXPECT().UpdateToken(user, "secure_token").Return(nil)

			authUser, err := accountService.Authorize("user@example.com", "123456")
			require.Equal(t, user, authUser)
			require.NoError(t, err)
		})
	})

}
