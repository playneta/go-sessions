package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/playneta/go-sessions/src/models"
	"github.com/stretchr/testify/require"
)

func TestRegister(t *testing.T) {
	request := []byte(`{"email": "user@example.com", "password": "123456"}`)

	ts, err := time.Parse(time.RFC1123, time.RFC1123)
	require.NoError(t, err)

	user := &models.User{
		ID:        1,
		Email:     "user@example.com",
		Password:  "password_hash",
		Token:     "",
		CreatedAt: ts,
		UpdatedAt: ts,
	}

	t.Run("Bad request", func(t *testing.T) {
		buf := bytes.NewBuffer([]byte(`i am not a good json`))

		suite := newTestSuite(t, http.MethodGet, buf, nil)
		defer suite.close()

		err := suite.api.Register(suite.context)
		require.Error(t, err)
	})

	t.Run("Error registering in serivce", func(t *testing.T) {
		buf := bytes.NewBuffer(request)
		suite := newTestSuite(t, http.MethodGet, buf, nil)
		defer suite.close()

		suite.accountService.EXPECT().Register("user@example.com", "123456").Return(nil, errors.New("service error!"))

		err := suite.api.Register(suite.context)
		require.Error(t, err, errors.New("service error!"))
	})

	t.Run("Success", func(t *testing.T) {
		buf := bytes.NewBuffer(request)
		suite := newTestSuite(t, http.MethodGet, buf, nil)
		defer suite.close()

		suite.accountService.EXPECT().Register("user@example.com", "123456").Return(user, nil)

		err := suite.api.Register(suite.context)
		require.NoError(t, err)
		{
			u := new(models.User)
			err := json.NewDecoder(suite.recorder.Body).Decode(u)
			require.NoError(t, err)
			require.Equal(t, u, user)
		}
	})
}
