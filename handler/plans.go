package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sknutsen/planner/database"
	"github.com/sknutsen/planner/lib"
	"github.com/sknutsen/planner/models"
	"github.com/sknutsen/planner/view"
)

func (h *Handler) Plan(c echo.Context) error {
	return h.EditPlan(c)
}

func (h *Handler) ListPlans(ctx context.Context, userId string) []database.Plan {
	dq := database.New(h.DB)

	plans, err := dq.ListPlans(ctx, userId)
	if err != nil {
		println(err.Error())
		return []database.Plan{}
	}

	return plans
}

func (h *Handler) DeletePlan(c echo.Context) error {
	return c.Redirect(http.StatusPermanentRedirect, fmt.Sprintf("/day/%s", lib.DateToString(time.Now())))
}

func (h *Handler) UpdatePlan(c echo.Context) error {
	var request models.UpdatePlanRequest

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

	if request.Id == "" {
		_, err = dq.CreatePlan(ctx, database.CreatePlanParams{
			Name: request.Name,
			User: user.UserId,
		})

		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed creating plan. err: %s", err))
		}
	} else {
		id, err := strconv.Atoi(request.Id)
		if err != nil {
			return c.String(http.StatusBadRequest, fmt.Sprintf("id is not a number. err: %s", err))
		}

		err = dq.UpdatePlan(ctx, database.UpdatePlanParams{
			ID:     int64(id),
			Name:   request.Name,
			UserId: user.UserId,
		})

		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed updating plan. err: %s", err))
		}
	}

	return h.Modal(c)
}

func (h *Handler) EditPlan(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.ErrBadRequest
	}

	planId, err := strconv.Atoi(id)
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

	plan, err := dq.GetPlan(ctx, database.GetPlanParams{
		ID:     int64(planId),
		UserId: state.UserProfile.UserId,
	})

	if err != nil {
		return err
	}

	component := view.Plan(state, models.Plan{
		Id:   int(plan.ID),
		Name: plan.Name,
	})
	return render(c, component)
}

func (h *Handler) CreatePlan(c echo.Context) error {
	state, err := models.GetClientState()
	if err != nil {
		println(err)
	}

	user, err := userProfileFromContext(c)
	if err != nil {
		println(err.Error())
	}

	state.UserProfile = user

	component := view.Plan(state, models.Plan{
		Id:   0,
		Name: "",
	})
	return render(c, component)
}
