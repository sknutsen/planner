package handler

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sknutsen/planner/lib"
	"github.com/sknutsen/planner/models"
	"github.com/sknutsen/planner/view"
)

func (h *Handler) Day(c echo.Context) error {
	date := c.Param("date")
	if date == "" {
		return echo.ErrBadRequest
	}

	state, err := models.GetClientState()
	if err != nil {
		println(err)
	}

	d, _ := lib.StringToDate(date)

	week := lib.ISOWeek(d)

	dates := lib.DatesInWeek(d.ISOWeek())
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

	weekday, err := state.Week.GetWeekday(d)
	if err != nil {
		return echo.ErrInternalServerError
	}

	component := view.Day(state, weekday)
	return component.Render(context.Background(), c.Response().Writer)
}

func (h *Handler) DayTasks(c echo.Context) error {
	date := c.Param("date")
	if date == "" {
		return echo.ErrBadRequest
	}

	d, _ := lib.StringToDate(date)

	component := view.DayTasks([]models.Task{
		{
			Id:          0,
			Date:        d,
			Title:       date,
			Subtitle:    "subtitle",
			Description: "description",
		},
		{
			Id:          0,
			Date:        d,
			Title:       date,
			Subtitle:    "subtitle",
			Description: "description",
		},
		{
			Id:          0,
			Date:        d,
			Title:       date,
			Subtitle:    "subtitle",
			Description: "description",
		},
		{
			Id:          0,
			Date:        d,
			Title:       date,
			Subtitle:    "subtitle",
			Description: "description",
		},
		{
			Id:          0,
			Date:        d,
			Title:       date,
			Subtitle:    "subtitle",
			Description: "description",
		},
		{
			Id:          0,
			Date:        d,
			Title:       date,
			Subtitle:    "subtitle",
			Description: "description",
		},
		{
			Id:          0,
			Date:        d,
			Title:       date,
			Subtitle:    "subtitle",
			Description: "description",
		},
		{
			Id:          0,
			Date:        d,
			Title:       date,
			Subtitle:    "subtitle",
			Description: "description",
		},
		{
			Id:          0,
			Date:        d,
			Title:       date,
			Subtitle:    "subtitle",
			Description: "description",
		},
		{
			Id:          0,
			Date:        d,
			Title:       date,
			Subtitle:    "subtitle",
			Description: "description",
		},
		{
			Id:          0,
			Date:        d,
			Title:       date,
			Subtitle:    "subtitle",
			Description: "description",
		},
		{
			Id:          0,
			Date:        d,
			Title:       date,
			Subtitle:    "subtitle",
			Description: "description",
		},
	})
	return component.Render(context.Background(), c.Response().Writer)
}
