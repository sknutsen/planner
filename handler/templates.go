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

func (h *Handler) Templates(c echo.Context) error {
	var planId int
	var err error
	id := c.Param("planId")
	if id != "" {
		planId, err = strconv.Atoi(id)
		if err != nil {
			println(err)
		}
	}

	state, err := models.GetTemplatesState()
	if err != nil {
		println(err)
	}

	user, err := userProfileFromContext(c)
	if err != nil {
		println(err.Error())
	}

	state.State.BaseRoute = routes.Templates

	state.State.UserProfile = user

	state.State.Plans = h.ListPlans(c.Request().Context(), state.State.UserProfile.UserId)

	state.State.SelectedPlanId = planid.Selected(state.State.Plans, planId)

	component := view.Templates(state)
	return render(c, component)
}

func (h *Handler) ListAllTemplates(c echo.Context) error {
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

	templates, err := dq.GetTemplatesByPlan(ctx, database.GetTemplatesByPlanParams{
		PlanId: int64(planId),
		UserId: user.UserId,
	})

	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed listing templates. err: %s", err))
	}

	component := view.PlanTemplates(models.PlanTemplatesResponse{
		Templates: models.TemplatesFromDBModels(templates),
	})
	return render(c, component)
}

func (h *Handler) Template(c echo.Context) error {
	return h.EditTemplate(c)
}

func (h *Handler) EditTemplate(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.ErrBadRequest
	}

	templateId, err := strconv.Atoi(id)
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

	template, err := dq.GetTemplate(ctx, database.GetTemplateParams{
		ID:     int64(templateId),
		UserId: state.UserProfile.UserId,
	})

	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed getting template. err: %s", err))
	}

	component := view.Template(state, models.Template{
		Id:          int(template.ID),
		Title:       template.Title,
		Subtitle:    lib.AsString(template.Subtitle),
		Description: lib.AsString(template.Description),
	})
	return render(c, component)
}

func (h *Handler) DeleteTemplate(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.ErrBadRequest
	}

	templateId, err := strconv.Atoi(id)
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

	template, err := dq.GetTemplate(ctx, database.GetTemplateParams{
		ID:     int64(templateId),
		UserId: state.UserProfile.UserId,
	})

	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed getting template. err: %s", err))
	}

	if err := dq.DeleteTemplate(ctx, database.DeleteTemplateParams{
		ID:     template.ID,
		UserId: state.UserProfile.UserId,
	}); err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed deleting template. err: %s", err))
	}

	c.Response().Header().Add("HX-Trigger", "updatedTemplate")

	return h.Modal(c)
}

func (h *Handler) CreateTemplate(c echo.Context) error {
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

	component := view.Template(state, models.Template{
		Id:          0,
		Title:       "",
		Subtitle:    "",
		Description: "",
	})
	return render(c, component)
}

func (h *Handler) UpdateTemplate(c echo.Context) error {
	var request models.UpdateTemplateRequest

	err := c.Bind(&request)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("bad request. err: %s", err))
	}

	return h.updateTemplate(c, request)
}

func (h *Handler) TemplateFromTask(c echo.Context) error {
	var request models.UpdateTemplateRequest

	err := c.Bind(&request)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("bad request. err: %s", err))
	}

	request.Id = "0"

	return h.updateTemplate(c, request)
}

func (h *Handler) updateTemplate(c echo.Context, r models.UpdateTemplateRequest) error {
	user, err := userProfileFromContext(c)
	if err != nil {
		println(err.Error())
	}

	ctx := c.Request().Context()
	dq := database.New(h.DB)

	if r.Id == "0" {
		println("Creating template")
		planId, err := strconv.Atoi(r.PlanId)
		if err != nil {
			return c.String(http.StatusBadRequest, fmt.Sprintf("id is not a number. err: %s", err))
		}

		_, err = dq.CreateTemplate(ctx, database.CreateTemplateParams{
			PlanID:      int64(planId),
			Title:       r.Title,
			Subtitle:    r.Subtitle,
			Description: r.Description,
		})

		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed creating template. err: %s", err))
		}
	} else {
		println("Updating template")
		id, err := strconv.Atoi(r.Id)
		if err != nil {
			return c.String(http.StatusBadRequest, fmt.Sprintf("id is not a number. err: %s", err))
		}

		err = dq.UpdateTemplate(ctx, database.UpdateTemplateParams{
			ID:          int64(id),
			Title:       r.Title,
			Subtitle:    r.Subtitle,
			Description: r.Description,
			UserId:      user.UserId,
		})

		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed updating template. err: %s", err))
		}
	}

	c.Response().Header().Add("HX-Trigger", "updatedTemplate")

	return h.Modal(c)
}
