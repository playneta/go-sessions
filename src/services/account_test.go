package services

import (
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/playneta/go-sessions/src/models"
	mock_providers "github.com/playneta/go-sessions/src/providers/mocks"
	mock_repositories "github.com/playneta/go-sessions/src/repositories/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestAccountService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Creating new account repository mock
	account := mock_repositories.NewMockUser(ctrl)
	messages := mock_repositories.NewMockMessage(ctrl)

	// Creating new hasher mock
	hasher := mock_providers.NewMockHasher(ctrl)

	// Noop logger
	logger := zap.NewNop().Sugar()

	// Service
	accountService := NewAccount(AccountOptions{
		AccountRepo: account,
		MessageRepo: messages,
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
				account.EXPECT().UpdateToken(user, gomock.Any()).Return(errors.New("token error!"))

				user, err := accountService.Authorize("user@example.com", "123456")
				require.Nil(t, user)
				require.Error(t, err, errors.New("token error!"))
			})
		})

		t.Run("Success", func(t *testing.T) {
			account.EXPECT().FindByEmail("user@example.com").Return(user, nil)
			hasher.EXPECT().Compare("123456", "my_password_hash").Return(true)
			account.EXPECT().UpdateToken(user, gomock.Any()).Return(nil)

			authUser, err := accountService.Authorize("user@example.com", "123456")
			require.Equal(t, user, authUser)
			require.NoError(t, err)
		})
	})

	t.Run("Create Message", func(t *testing.T) {
		user := models.User{
			ID:        1,
			Email:     "user@example.com",
			Password:  "my_password_hash",
			Token:     "",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		t.Run("Empty text", func(t *testing.T) {
			message, err := accountService.CreateMessage(user, "", "")
			require.Nil(t, message)
			require.Error(t, err)
		})

		t.Run("Non existent receiver", func(t *testing.T) {
			account.EXPECT().FindByEmail("unknown@example.com").Return(nil, errors.New("unknown user"))

			message, err := accountService.CreateMessage(user, "unknown@example.com", "text")
			require.Nil(t, message)
			require.Error(t, err)
		})

		t.Run("Failure", func(t *testing.T) {
			messages.EXPECT().Create(gomock.Any()).Return(errors.New("error creating message"))

			message, err := accountService.CreateMessage(user, "", "text")
			require.Nil(t, message)
			require.Error(t, err)
		})

		t.Run("Success", func(t *testing.T) {
			account.EXPECT().FindByEmail("user@example.com").Return(&models.User{ID: 100}, nil)
			messages.EXPECT().Create(gomock.Any()).Return(nil)

			message, err := accountService.CreateMessage(user, "user@example.com", "text")

			require.NoError(t, err)
			require.Equal(t, user.ID, message.UserId)
			require.Equal(t, int64(100), message.ReceiverId)
			require.Equal(t, "text", message.Text)
		})
	})

	t.Run("History", func(t *testing.T) {
		user := models.User{
			ID:        1,
			Email:     "user@example.com",
			Password:  "my_password_hash",
			Token:     "",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		t.Run("Errors", func(t *testing.T) {
			t.Run("Public", func(t *testing.T) {
				messages.EXPECT().LastPublicMessages(10).Return(nil, errors.New("error getting public messages"))

				messages, err := accountService.History(user)
				require.Nil(t, messages)
				require.Error(t, err)
			})

			t.Run("Private", func(t *testing.T) {
				messages.EXPECT().LastPublicMessages(10).Return([]models.Message{}, nil)
				messages.EXPECT().LastPrivateMessages(user, 10).Return(nil, errors.New("error getting private messages"))

				messages, err := accountService.History(user)
				require.Nil(t, messages)
				require.Error(t, err)
			})
		})

		t.Run("Success", func(t *testing.T) {
			ts := time.Now()

			public := []models.Message{
				{
					Id:        1,
					CreatedAt: ts.Add(-1 * time.Minute),
				},
				{
					Id:        2,
					CreatedAt: ts,
				},
			}

			private := []models.Message{
				{
					Id:        3,
					CreatedAt: ts,
				},
			}

			messages.EXPECT().LastPublicMessages(10).Return(public, nil)
			messages.EXPECT().LastPrivateMessages(user, 10).Return(private, nil)

			messages, err := accountService.History(user)
			require.NoError(t, err)
			require.Len(t, messages, 3)
		})
	})

}
