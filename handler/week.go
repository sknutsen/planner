package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/sknutsen/planner/lib"
	"github.com/sknutsen/planner/models"
	"github.com/sknutsen/planner/view"
)

func (h *Handler) Index(c echo.Context) error {
	week := c.Param("week")
	if week == "" {
		week = lib.ISOWeek(time.Now())
	}

	state, err := models.GetClientState()
	if err != nil {
		println(err)
	}

	sess, err := session.Get("session", c)
	if err != nil {
		println(err)
	}

	state.UserProfile = models.GetUserProfile(sess.Values["profile"].(map[string]interface{}))

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

	state.Plans = h.ListPlans(state.UserProfile.UserId)

	if len(state.Plans) > 0 {
		state.SelectedPlanId = int(state.Plans[0].ID)
	}

	component := view.Index(state)
	return component.Render(context.Background(), c.Response().Writer)
}
