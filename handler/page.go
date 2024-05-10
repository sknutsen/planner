package handler

import (
	"context"
	"time"

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

	component := view.Index(state)
	return component.Render(context.Background(), c.Response().Writer)
}
