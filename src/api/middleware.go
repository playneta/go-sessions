package api

import (
	"net/http"

	"github.com/labstack/echo"
)

func (a *API) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("X-TOKEN")

		user, err := a.userRepo.FindByToken(token)
		if err != nil || user == nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "empty/wrong token")
		}

		c.Set("user", user)
		return next(c)
	}
}
