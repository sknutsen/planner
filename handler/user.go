package handler

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/sknutsen/planner/models"
	"github.com/sknutsen/planner/view"
)

func (h *Handler) User(c echo.Context) error {
	state, err := models.GetUserState()
	if err != nil {
		return err
	}

	user, err := userProfileFromContext(c)
	if err != nil {
		println(err.Error())
	}

	state.State.UserProfile = user

	component := view.User(state)
	return component.Render(context.Background(), c.Response().Writer)
}
