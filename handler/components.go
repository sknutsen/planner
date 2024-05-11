package handler

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/sknutsen/planner/view"
)

func (h *Handler) Modal(c echo.Context) error {
	component := view.Modal("hidden")
	return component.Render(context.Background(), c.Response().Writer)
}
