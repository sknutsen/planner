package handler

import (
	"context"
	"strconv"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/sknutsen/planner/models"
	"github.com/sknutsen/planner/routes"
	"github.com/sknutsen/planner/view"
)

func (h *Handler) History(c echo.Context) error {
	var planId int
	var err error
	id := c.Param("planId")
	if id != "" {
		planId, err = strconv.Atoi(id)
		if err != nil {
			println(err)
		}
	}

	state, err := models.GetHistoryState()
	if err != nil {
		println(err)
	}

	sess, err := session.Get("session", c)
	if err != nil {
		println(err)
	}

	state.State.BaseRoute = routes.History

	state.State.UserProfile = models.GetUserProfile(sess.Values["profile"].(map[string]interface{}))

	state.State.Plans = h.ListPlans(state.State.UserProfile.UserId)

	if len(state.State.Plans) > 0 {
		for _, p := range state.State.Plans {
			if planId == int(p.ID) || planId == 0 {
				state.State.SelectedPlanId = int(p.ID)
				break
			}
		}
	}

	component := view.History(state)
	return component.Render(context.Background(), c.Response().Writer)
}
