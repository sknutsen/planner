package handler

import (
	"context"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/sknutsen/planner/models"
	"github.com/sknutsen/planner/view"
)

func (h *Handler) User(c echo.Context) error {
	state, err := models.GetUserState()
	if err != nil {
		return err
	}

	sess, err := session.Get("session", c)
	if err != nil {
		println(err)
	}

	state.State.UserProfile = models.GetUserProfile(sess.Values["profile"].(map[string]interface{}))

	// for k, v := range sess.Values["profile"].(map[string]interface{}) {
	// 	fmt.Printf("k: %v v: %v\n", k, v)
	// }

	component := view.User(state)
	return component.Render(context.Background(), c.Response().Writer)
}
