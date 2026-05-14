package handler

import (
	"context"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sknutsen/planner/lib"
	"github.com/sknutsen/planner/models"
	"github.com/sknutsen/planner/routes"
	"github.com/sknutsen/planner/view"
)

func (h *Handler) Week(c echo.Context) error {
	var planId int
	id := c.Param("id")
	if id != "" {
		planId, _ = strconv.Atoi(id)
	}

	week := c.Param("week")
	if week == "" {
		week = lib.ISOWeek(time.Now())
	}

	state, err := models.GetWeekState()
	if err != nil {
		println(err)
	}

	user, err := userProfileFromContext(c)
	if err != nil {
		println(err.Error())
	}

	state.State.BaseRoute = routes.Week

	state.State.UserProfile = user

	dates := lib.DatesInWeek(lib.ISOWeekFromString(week))
	state.Week.ISOWeek = week

	models.PopulateWeekDates(&state.Week, dates)

	state.State.Plans = h.ListPlans(state.State.UserProfile.UserId)

	state.State.SelectedPlanId = selectedPlanID(state.State.Plans, planId)

	component := view.Index(state)
	return component.Render(context.Background(), c.Response().Writer)
}
