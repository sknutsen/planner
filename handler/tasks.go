package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sknutsen/planner/lib"
	"github.com/sknutsen/planner/models"
	"github.com/sknutsen/planner/view"
)

func (h *Handler) Task(c echo.Context) error {
	return h.EditTask(c)
}

func (h *Handler) EditTask(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.ErrBadRequest
	}

	taskId, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	state, err := models.GetClientState()
	if err != nil {
		println(err)
	}

	component := view.Task(state, models.Task{
		Id:          taskId,
		Date:        time.Now(),
		Title:       "Title",
		Subtitle:    "Subtitle",
		Description: "Description",
	})
	return component.Render(context.Background(), c.Response().Writer)
}

func (h *Handler) UpdateTask(c echo.Context) error {
	return h.Modal(c)
}

func (h *Handler) DeleteTask(c echo.Context) error {
	return c.Redirect(http.StatusPermanentRedirect, fmt.Sprintf("/day/%s", lib.DateToString(time.Now())))
}

func (h *Handler) CreateTask(c echo.Context) error {
	date := c.Param("date")
	if date == "" {
		return echo.ErrBadRequest
	}

	d, _ := lib.StringToDate(date)

	state, err := models.GetClientState()
	if err != nil {
		println(err)
	}

	component := view.Task(state, models.Task{
		Id:          0,
		Date:        d,
		Title:       "",
		Subtitle:    "",
		Description: "hello world",
	})
	return component.Render(context.Background(), c.Response().Writer)
}
