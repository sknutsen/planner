package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo-contrib/session"
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

	sess, err := session.Get("session", c)
	if err != nil {
		println(err)
	}

	state.State.SelectedPlanId = planId
	state.State.UserProfile = models.GetUserProfile(sess.Values["profile"].(map[string]interface{}))

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

	sess, err := session.Get("session", c)
	if err != nil {
		println(err)
	}

	user := models.GetUserProfile(sess.Values["profile"].(map[string]interface{}))

	db := h.openDB()
	defer db.Close()

	ctx := context.Background()
	dq := database.New(db)

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
