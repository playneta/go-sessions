package api

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/playneta/go-sessions/src/models"
)

func (a *Api) Profile(ctx echo.Context) error {
	user := ctx.Get("user").(*models.User)
	return ctx.JSON(http.StatusOK, user)
}
