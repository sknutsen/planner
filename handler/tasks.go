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

	user, err := userProfileFromContext(c)
	if err != nil {
		println(err.Error())
	}

	state.UserProfile = user

	ctx := c.Request().Context()
	dq := database.New(h.DB)

	task, err := dq.GetTask(ctx, database.GetTaskParams{
		ID:     int64(taskId),
		UserId: state.UserProfile.UserId,
	})

	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed getting task. err: %s", err))
	}

	println(task.Date)

	state.SelectedPlanId = int(task.PlanID)

	component := view.Task(state, models.Task{
		Id:          int(task.ID),
		Date:        lib.StripDateString(task.Date),
		Title:       task.Title,
		Subtitle:    lib.AsString(task.Subtitle),
		Description: lib.AsString(task.Description),
		IsComplete:  task.IsComplete != 0,
	})
	return component.Render(context.Background(), c.Response().Writer)
}

func (h *Handler) CopyTask(c echo.Context) error {
	var request models.UpdateTaskRequest

	err := c.Bind(&request)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("bad request. err: %s", err))
	}

	user, err := userProfileFromContext(c)
	if err != nil {
		println(err.Error())
	}

	println(user.UserId)
	if user.UserId != "" {

		ctx := c.Request().Context()
		dq := database.New(h.DB)

		taskId, err := strconv.Atoi(request.Id)
		if err != nil {
			return c.String(http.StatusBadRequest, fmt.Sprintf("id is not a number. err: %s", err))
		}

		task, err := dq.GetTask(ctx, database.GetTaskParams{
			ID:     int64(taskId),
			UserId: user.UserId,
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

	user, err := userProfileFromContext(c)
	if err != nil {
		println(err.Error())
	}

	ctx := c.Request().Context()
	dq := database.New(h.DB)

	if request.Id == "0" {
		planId, err := strconv.Atoi(request.PlanId)
		if err != nil {
			return c.String(http.StatusBadRequest, fmt.Sprintf("id is not a number. err: %s", err))
		}

		if request.Template != "" {
			templateId, err := strconv.Atoi(request.Template)
			if err != nil {
				return c.String(http.StatusBadRequest, fmt.Sprintf("id is not a number. err: %s", err))
			}
			err = dq.CreateTaskFromTemplate(ctx, database.CreateTaskFromTemplateParams{
				TemplateId: int64(templateId),
				Date:       request.Date,
				UserId:     user.UserId,
			})
		} else {
			err = dq.CreateTask(ctx, database.CreateTaskParams{
				PlanID:      int64(planId),
				Title:       request.Title,
				Subtitle:    request.Subtitle,
				Description: request.Description,
				Date:        request.Date,
			})
		}

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
			UserId:      user.UserId,
		})

		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed updating task. err: %s", err))
		}
	}

	c.Response().Header().Add("HX-Trigger", "updatedTask")

	return h.Modal(c)
}

func (h *Handler) ToggleIsCompleteTask(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.ErrBadRequest
	}

	taskId, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	user, err := userProfileFromContext(c)
	if err != nil {
		println(err.Error())
	}

	ctx := c.Request().Context()
	dq := database.New(h.DB)

	task, err := dq.GetTask(ctx, database.GetTaskParams{
		ID:     int64(taskId),
		UserId: user.UserId,
	})

	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed getting task. err: %s", err))
	}

	var isComplete int64
	if task.IsComplete == 0 {
		isComplete = 1
	} else {
		isComplete = 0
	}

	err = dq.SetIsCompleteTask(ctx, database.SetIsCompleteTaskParams{
		IsComplete: isComplete,
		ID:         int64(taskId),
		UserId:     user.UserId,
	})

	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed updating task completion. err: %s", err))
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

	user, err := userProfileFromContext(c)
	if err != nil {
		println(err.Error())
	}

	state.UserProfile = user

	ctx := c.Request().Context()
	dq := database.New(h.DB)

	task, err := dq.GetTask(ctx, database.GetTaskParams{
		ID:     int64(taskId),
		UserId: state.UserProfile.UserId,
	})

	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed getting task. err: %s", err))
	}

	if err := dq.DeleteTask(ctx, database.DeleteTaskParams{
		ID:     task.ID,
		UserId: state.UserProfile.UserId,
	}); err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed deleting task. err: %s", err))
	}

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

	// sess, err := session.Get("session", c)
	// if err != nil {
	// 	println(err)
	// }

	state.SelectedPlanId = planId
	user, err := userProfileFromContext(c)
	if err != nil {
		println(err.Error())
	}

	state.UserProfile = user

	// state.UserProfile = models.GetUserProfile(sess.Values["profile"].(map[string]interface{}))

	component := view.Task(state, models.Task{
		Id:          0,
		Date:        date,
		Title:       "",
		Subtitle:    "",
		Description: "",
	})
	return component.Render(context.Background(), c.Response().Writer)
}

func (h *Handler) CreateTaskFromTemplate(c echo.Context) error {
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

	// sess, err := session.Get("session", c)
	// if err != nil {
	// 	println(err)
	// }

	state.SelectedPlanId = planId
	// state.UserProfile = models.GetUserProfile(sess.Values["profile"].(map[string]interface{}))

	user, err := userProfileFromContext(c)
	if err != nil {
		println(err.Error())
	}

	state.UserProfile = user

	ctx := c.Request().Context()
	dq := database.New(h.DB)

	templates, err := dq.GetTemplatesByPlan(ctx, database.GetTemplatesByPlanParams{
		PlanId: int64(planId),
		UserId: state.UserProfile.UserId,
	})

	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed listing templates. err: %s", err))
	}

	component := view.TaskFromTemplate(state, models.Task{
		Id:          0,
		Date:        date,
		Title:       "",
		Subtitle:    "",
		Description: "",
	}, models.TemplatesFromDBModels(templates))
	return component.Render(context.Background(), c.Response().Writer)
}

func (h *Handler) ListAllTasks(c echo.Context) error {
	var planId int
	id := c.Param("planId")
	planId, err := strconv.Atoi(id)
	if err != nil {
		return echo.ErrBadRequest
	}

	user, err := userProfileFromContext(c)
	if err != nil {
		println(err.Error())
	}

	ctx := c.Request().Context()
	dq := database.New(h.DB)

	tasks, err := dq.GetTasksByPlan(ctx, database.GetTasksByPlanParams{
		PlanId: int64(planId),
		UserId: user.UserId,
	})

	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed listing tasks for plan. err: %s", err))
	}

	component := view.HistoryTasks(models.HistoryTasksResponse{
		Tasks: models.TasksFromDBModels(tasks),
	})
	return component.Render(context.Background(), c.Response().Writer)
}
