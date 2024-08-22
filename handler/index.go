package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sknutsen/planner/routes"
)

func (handler *Handler) Index(c echo.Context) error {
	return c.Redirect(http.StatusTemporaryRedirect, routes.Week)
}
