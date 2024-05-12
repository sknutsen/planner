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

func (h *Handler) Plan(c echo.Context) error {
	return h.EditPlan(c)
}

func (h *Handler) ListPlans(userId string) []database.Plan {
	db := h.openDB()
	defer db.Close()

	ctx := context.Background()
	dq := database.New(db)

	plans, err := dq.ListPlans(ctx, database.ListPlansParams{
		User:   userId,
		User_2: userId,
	})
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
		err = dq.CreatePlan(ctx, database.CreatePlanParams{
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
			ID:   int64(id),
			Name: request.Name,
			User: user.UserId,
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

	component := view.Plan(state, models.Plan{
		Id:   taskId,
		Name: "",
	})
	return component.Render(context.Background(), c.Response().Writer)
}

func (h *Handler) CreatePlan(c echo.Context) error {
	state, err := models.GetClientState()
	if err != nil {
		println(err)
	}

	sess, err := session.Get("session", c)
	if err != nil {
		println(err)
	}

	state.UserProfile = models.GetUserProfile(sess.Values["profile"].(map[string]interface{}))

	component := view.Plan(state, models.Plan{
		Id:   0,
		Name: "",
	})
	return component.Render(context.Background(), c.Response().Writer)
}
