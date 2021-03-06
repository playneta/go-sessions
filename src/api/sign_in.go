package api

import (
	"net/http"

	"github.com/labstack/echo"
)

func (a *API) SignIn(ctx echo.Context) error {
	var req UserRequest
	if err := ctx.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := a.accountService.Authorize(req.Email, req.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	return ctx.JSON(200, user)
}
