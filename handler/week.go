package handler

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/labstack/echo-contrib/session"
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

	sess, err := session.Get("session", c)
	if err != nil {
		println(err)
	}

	state.State.BaseRoute = routes.Week

	state.State.UserProfile = models.GetUserProfile(sess.Values["profile"].(map[string]interface{}))

	fmt.Printf("week: %s\n", week)
	dates := lib.DatesInWeek(lib.ISOWeekFromString(week))
	state.Week.ISOWeek = week

	for _, d := range dates {
		date, _ := lib.StringToDate(d)

		switch date.Weekday() {
		case time.Monday:
			state.Week.Monday.Date = date
		case time.Tuesday:
			state.Week.Tuesday.Date = date
		case time.Wednesday:
			state.Week.Wednesday.Date = date
		case time.Thursday:
			state.Week.Thursday.Date = date
		case time.Friday:
			state.Week.Friday.Date = date
		case time.Saturday:
			state.Week.Saturday.Date = date
		case time.Sunday:
			state.Week.Sunday.Date = date
		}
	}

	state.State.Plans = h.ListPlans(state.State.UserProfile.UserId)

	if len(state.State.Plans) > 0 {
		for _, p := range state.State.Plans {
			if planId == int(p.ID) || planId == 0 {
				state.State.SelectedPlanId = int(p.ID)
				break
			}
		}
	}

	component := view.Index(state)
	return component.Render(context.Background(), c.Response().Writer)
}
