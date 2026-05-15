package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sknutsen/planner/database"
	"github.com/sknutsen/planner/internal/planid"
	"github.com/sknutsen/planner/lib"
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

	user, err := userProfileFromContext(c)
	if err != nil {
		println(err.Error())
	}

	state.State.BaseRoute = routes.History

	state.State.UserProfile = user

	state.State.Plans = h.ListPlans(c.Request().Context(), state.State.UserProfile.UserId)

	state.State.SelectedPlanId = planid.Selected(state.State.Plans, planId)

	component := view.Resources(state)
	return render(c, component)
}

func (h *Handler) ListAllResources(c echo.Context) error {
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

	tasks, err := dq.GetResourcesByPlan(ctx, database.GetResourcesByPlanParams{
		PlanId: int64(planId),
		UserId: user.UserId,
	})

	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed listing resources. err: %s", err))
	}

	component := view.PlanResources(models.PlanResourcesResponse{
		Resources: models.ResourcesFromDBModels(tasks),
	})
	return render(c, component)
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

	user, err := userProfileFromContext(c)
	if err != nil {
		println(err.Error())
	}

	state.UserProfile = user

	ctx := c.Request().Context()
	dq := database.New(h.DB)

	task, err := dq.GetResource(ctx, database.GetResourceParams{
		ID:     int64(taskId),
		UserId: state.UserProfile.UserId,
	})

	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed getting resource. err: %s", err))
	}

	component := view.Resource(state, models.Resource{
		Id:      int(task.ID),
		Title:   task.Title,
		Type:    int(task.ResourceType),
		Content: lib.AsString(task.Content),
	})
	return render(c, component)
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

	user, err := userProfileFromContext(c)
	if err != nil {
		println(err.Error())
	}

	state.UserProfile = user

	ctx := c.Request().Context()
	dq := database.New(h.DB)

	task, err := dq.GetResource(ctx, database.GetResourceParams{
		ID:     int64(taskId),
		UserId: state.UserProfile.UserId,
	})

	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed getting resource. err: %s", err))
	}

	if err := dq.DeleteResource(ctx, database.DeleteResourceParams{
		ID:     task.ID,
		UserId: state.UserProfile.UserId,
	}); err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed deleting resource. err: %s", err))
	}

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

	user, err := userProfileFromContext(c)
	if err != nil {
		println(err.Error())
	}

	state.SelectedPlanId = planId
	state.UserProfile = user

	component := view.Resource(state, models.Resource{
		Id:      0,
		Title:   "",
		Type:    0,
		Content: "",
	})
	return render(c, component)
}

func (h *Handler) UpdateResource(c echo.Context) error {
	var request models.UpdateResourceRequest

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
			UserId:       user.UserId,
		})

		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed updating resource. err: %s", err))
		}
	}

	c.Response().Header().Add("HX-Trigger", "updatedResource")

	return h.Modal(c)
}
