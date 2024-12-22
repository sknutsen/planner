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
	"github.com/sknutsen/planner/routes"
	"github.com/sknutsen/planner/view"
)

func (h *Handler) Resources(c echo.Context) error {
	var planId int
	var err error
	id := c.Param("planId")
	if id != "" {
		planId, err = strconv.Atoi(id)
		if err != nil {
			println(err)
		}
	}

	state, err := models.GetResourcesState()
	if err != nil {
		println(err)
	}

	sess, err := session.Get("session", c)
	if err != nil {
		println(err)
	}

	state.State.BaseRoute = routes.History

	state.State.UserProfile = models.GetUserProfile(sess.Values["profile"].(map[string]interface{}))

	state.State.Plans = h.ListPlans(state.State.UserProfile.UserId)

	if len(state.State.Plans) > 0 {
		for _, p := range state.State.Plans {
			if planId == int(p.ID) || planId == 0 {
				state.State.SelectedPlanId = int(p.ID)
				break
			}
		}
	}

	component := view.Resources(state)
	return component.Render(context.Background(), c.Response().Writer)
}

func (h *Handler) ListAllResources(c echo.Context) error {
	var planId int
	id := c.Param("planId")
	planId, err := strconv.Atoi(id)
	if err != nil {
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

	tasks, err := dq.GetResourcesByPlan(ctx, database.GetResourcesByPlanParams{
		ID:     int64(planId),
		User:   user.UserId,
		User_2: user.UserId,
	})

	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed listing tasks by date. err: %s", err))
	}

	component := view.PlanResources(models.PlanResourcesResponse{
		Resources: models.ResourcesFromDBModels(tasks),
	})
	return component.Render(context.Background(), c.Response().Writer)
}

func (h *Handler) Resource(c echo.Context) error {
	return h.EditResource(c)
}

func (h *Handler) EditResource(c echo.Context) error {
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

	task, err := dq.GetResource(ctx, database.GetResourceParams{
		ID:     int64(taskId),
		User:   state.UserProfile.UserId,
		User_2: state.UserProfile.UserId,
	})

	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed getting resource. err: %s", err))
	}

	component := view.Resource(state, models.Resource{
		Id:      int(task.ID),
		Title:   task.Title,
		Type:    int(task.ResourceType),
		Content: task.Content.(string),
	})
	return component.Render(context.Background(), c.Response().Writer)
}

func (h *Handler) DeleteResource(c echo.Context) error {
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

	task, err := dq.GetResource(ctx, database.GetResourceParams{
		ID:     int64(taskId),
		User:   state.UserProfile.UserId,
		User_2: state.UserProfile.UserId,
	})

	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed getting resource. err: %s", err))
	}

	dq.DeleteResource(ctx, database.DeleteResourceParams{
		ID:     task.ID,
		User:   state.UserProfile.UserId,
		User_2: state.UserProfile.UserId,
	})

	c.Response().Header().Add("HX-Trigger", "updatedResource")

	return h.Modal(c)
}

func (h *Handler) CreateResource(c echo.Context) error {
	var planId int
	id := c.Param("planId")
	planId, err := strconv.Atoi(id)
	if err != nil {
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

	component := view.Resource(state, models.Resource{
		Id:      0,
		Title:   "",
		Type:    0,
		Content: "",
	})
	return component.Render(context.Background(), c.Response().Writer)
}

func (h *Handler) UpdateResource(c echo.Context) error {
	var request models.UpdateResourceRequest

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

	resourceType, err := strconv.Atoi(request.Type)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("type is not a number. err: %s", err))
	}

	if request.Id == "0" {
		planId, err := strconv.Atoi(request.PlanId)
		if err != nil {
			return c.String(http.StatusBadRequest, fmt.Sprintf("id is not a number. err: %s", err))
		}

		err = dq.CreateResource(ctx, database.CreateResourceParams{
			PlanID:       int64(planId),
			Title:        request.Title,
			ResourceType: int64(resourceType),
			Content:      request.Content,
		})

		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed creating resource. err: %s", err))
		}
	} else {
		id, err := strconv.Atoi(request.Id)
		if err != nil {
			return c.String(http.StatusBadRequest, fmt.Sprintf("id is not a number. err: %s", err))
		}

		err = dq.UpdateResource(ctx, database.UpdateResourceParams{
			ID:           int64(id),
			Title:        request.Title,
			ResourceType: int64(resourceType),
			Content:      request.Content,
			User:         user.UserId,
			User_2:       user.UserId,
		})

		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed updating resource. err: %s", err))
		}
	}

	c.Response().Header().Add("HX-Trigger", "updatedResource")

	return h.Modal(c)
}
