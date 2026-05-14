package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sknutsen/planner/database"
	"github.com/sknutsen/planner/lib"
	"github.com/sknutsen/planner/models"
	"github.com/sknutsen/planner/view"
)

func (h *Handler) Day(c echo.Context) error {
	var planId int
	id := c.Param("planId")
	planId, err := strconv.Atoi(id)
	if err != nil {
		return echo.ErrBadRequest
	}

	date := c.Param("date")
	if date == "" {
		return echo.ErrBadRequest
	}

	state, err := models.GetWeekState()
	if err != nil {
		println(err)
	}

	d, _ := lib.StringToDate(date)

	week := lib.ISOWeek(d)

	dates := lib.DatesInWeek(d.ISOWeek())
	state.Week.ISOWeek = week

	user, err := userProfileFromContext(c)
	if err != nil {
		println(err.Error())
	}

	state.State.SelectedPlanId = planId
	state.State.UserProfile = user

	models.PopulateWeekDates(&state.Week, dates)

	weekday, err := state.Week.GetWeekday(d)
	if err != nil {
		return echo.ErrInternalServerError
	}

	component := view.Day(state, weekday)
	return component.Render(context.Background(), c.Response().Writer)
}

func (h *Handler) DayTasks(c echo.Context) error {
	var planId int
	id := c.Param("planId")
	planId, err := strconv.Atoi(id)
	if err != nil {
		return echo.ErrBadRequest
	}

	date := c.Param("date")
	if date == "" {
		return echo.ErrBadRequest
	}

	user, err := userProfileFromContext(c)
	if err != nil {
		println(err.Error())
	}

	ctx := c.Request().Context()
	dq := database.New(h.DB)

	tasks, err := dq.GetTasksByDate(ctx, database.GetTasksByDateParams{
		Date:   date,
		PlanId: int64(planId),
		UserId: user.UserId,
	})

	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed listing tasks by date. err: %s", err))
	}

	component := view.DayTasks(models.DayTasksResponse{
		Date:            date,
		Tasks:           models.TasksFromDBModels(tasks),
		HideDescription: true,
	})
	return component.Render(context.Background(), c.Response().Writer)
}
