package handler

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sknutsen/planner/internal/planid"
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

	user, err := userProfileFromContext(c)
	if err != nil {
		println(err.Error())
	}

	state.State.BaseRoute = routes.History

	state.State.UserProfile = user

	state.State.Plans = h.ListPlans(c.Request().Context(), state.State.UserProfile.UserId)

	state.State.SelectedPlanId = planid.Selected(state.State.Plans, planId)

	return render(c, view.History(state))
}
