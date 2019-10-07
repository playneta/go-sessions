package api

import (
	"net/http"
	"testing"

	"github.com/labstack/echo"
	"github.com/playneta/go-sessions/src/models"
	"github.com/stretchr/testify/require"
)

func TestMiddleware(t *testing.T) {
	t.Run("Bad token", func(t *testing.T) {
		suite := newTestSuite(t, http.MethodGet, nil, nil)
		defer suite.close()

		suite.userRepo.EXPECT().FindByToken("").Return(nil, nil)

		next := func(c echo.Context) error {
			return nil
		}
		handler := suite.api.AuthMiddleware(next)
		require.NotNil(t, handler)
		err := handler(suite.context)
		require.Error(t, err, echo.NewHTTPError(http.StatusUnauthorized, "empty/wrong token"))
	})

	t.Run("Good token", func(t *testing.T) {
		suite := newTestSuite(t, http.MethodGet, nil, map[string]string{
			"X-TOKEN": "token",
		})
		defer suite.close()

		user := &models.User{
			ID:    1,
			Email: "user@example.com",
		}

		suite.userRepo.EXPECT().FindByToken("token").Return(user, nil)

		next := func(c echo.Context) error {
			return nil
		}
		handler := suite.api.AuthMiddleware(next)
		require.NotNil(t, handler)
		err := handler(suite.context)
		require.NoError(t, err)
		require.Nil(t, err)
	})
}
