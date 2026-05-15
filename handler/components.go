package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/sknutsen/planner/view"
)

func (h *Handler) Modal(c echo.Context) error {
	return render(c, view.Modal("hidden"))
}
