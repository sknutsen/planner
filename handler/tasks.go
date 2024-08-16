package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/sknutsen/planner/database"
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

	sess, err := session.Get("session", c)
	if err != nil {
		println(err)
	}

	state.UserProfile = models.GetUserProfile(sess.Values["profile"].(map[string]interface{}))

	db := h.openDB()
	defer db.Close()

	ctx := context.Background()
	dq := database.New(db)

	task, err := dq.GetTask(ctx, database.GetTaskParams{
		ID:     int64(taskId),
		User:   state.UserProfile.UserId,
		User_2: state.UserProfile.UserId,
	})

	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed getting task. err: %s", err))
	}

	component := view.Task(state, models.Task{
		Id:          int(task.ID),
		Date:        task.Date,
		Title:       task.Title,
		Subtitle:    task.Subtitle.(string),
		Description: task.Description.(string),
	})
	return component.Render(context.Background(), c.Response().Writer)
}

func (h *Handler) CopyTask(c echo.Context) error {
	var request models.UpdateTaskRequest

	err := c.Bind(&request)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("bad request. err: %s", err))
	}

	sess, err := session.Get("session", c)
	if err != nil {
		println(err)
	}

	user := models.GetUserProfile(sess.Values["profile"].(map[string]interface{}))
	println(user.UserId)
	if user.UserId != "" {

		db := h.openDB()
		defer db.Close()

		ctx := context.Background()
		dq := database.New(db)

		taskId, err := strconv.Atoi(request.Id)
		if err != nil {
			return c.String(http.StatusBadRequest, fmt.Sprintf("id is not a number. err: %s", err))
		}

		task, err := dq.GetTask(ctx, database.GetTaskParams{
			ID:     int64(taskId),
			User:   user.UserId,
			User_2: user.UserId,
		})

		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed getting task. err: %s", err))
		}

		err = dq.CreateTask(ctx, database.CreateTaskParams{
			PlanID:      int64(task.PlanID),
			Title:       request.Title,
			Subtitle:    request.Subtitle,
			Description: request.Description,
			Date:        request.Date,
		})

		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed creating task. err: %s", err))
		}
	}

	c.Response().Header().Add("HX-Trigger", "updatedTask")

	return h.Modal(c)
}

func (h *Handler) UpdateTask(c echo.Context) error {
	var request models.UpdateTaskRequest

	err := c.Bind(&request)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("bad request. err: %s", err))
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

	if request.Id == "0" {
		planId, err := strconv.Atoi(request.PlanId)
		if err != nil {
			return c.String(http.StatusBadRequest, fmt.Sprintf("id is not a number. err: %s", err))
		}

		err = dq.CreateTask(ctx, database.CreateTaskParams{
			PlanID:      int64(planId),
			Title:       request.Title,
			Subtitle:    request.Subtitle,
			Description: request.Description,
			Date:        request.Date,
		})

		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed creating task. err: %s", err))
		}
	} else {
		id, err := strconv.Atoi(request.Id)
		if err != nil {
			return c.String(http.StatusBadRequest, fmt.Sprintf("id is not a number. err: %s", err))
		}

		err = dq.UpdateTask(ctx, database.UpdateTaskParams{
			ID:          int64(id),
			Title:       request.Title,
			Subtitle:    request.Subtitle,
			Description: request.Description,
			Date:        request.Date,
			User:        user.UserId,
			User_2:      user.UserId,
		})

		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed updating task. err: %s", err))
		}
	}

	c.Response().Header().Add("HX-Trigger", "updatedTask")

	return h.Modal(c)
}

func (h *Handler) DeleteTask(c echo.Context) error {
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

	sess, err := session.Get("session", c)
	if err != nil {
		println(err)
	}

	state.UserProfile = models.GetUserProfile(sess.Values["profile"].(map[string]interface{}))

	db := h.openDB()
	defer db.Close()

	ctx := context.Background()
	dq := database.New(db)

	task, err := dq.GetTask(ctx, database.GetTaskParams{
		ID:     int64(taskId),
		User:   state.UserProfile.UserId,
		User_2: state.UserProfile.UserId,
	})

	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed getting task. err: %s", err))
	}

	dq.DeleteTask(ctx, database.DeleteTaskParams{
		ID:     task.ID,
		User:   state.UserProfile.UserId,
		User_2: state.UserProfile.UserId,
	})

	c.Response().Header().Add("HX-Trigger", "updatedTask")

	return h.Modal(c)
}

func (h *Handler) CreateTask(c echo.Context) error {
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

	state, err := models.GetClientState()
	if err != nil {
		println(err)
	}

	sess, err := session.Get("session", c)
	if err != nil {
		println(err)
	}

	state.SelectedPlanId = planId
	state.UserProfile = models.GetUserProfile(sess.Values["profile"].(map[string]interface{}))

	component := view.Task(state, models.Task{
		Id:          0,
		Date:        date,
		Title:       "",
		Subtitle:    "",
		Description: "",
	})
	return component.Render(context.Background(), c.Response().Writer)
}
