package handler

// Mobile JSON API (/api/v1). Optimistic concurrency: send If-Match: W/"<updated_at>" or JSON
// base_updated_at on mutating requests; on mismatch the server responds 409 with error.entity
// containing the current row.

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/sknutsen/planner/database"
	"github.com/sknutsen/planner/internal/apijson"
	"github.com/sknutsen/planner/internal/synccursor"
)

func parseLimit(s string, def, max int) int {
	if s == "" {
		return def
	}
	n, err := strconv.Atoi(s)
	if err != nil || n < 1 {
		return def
	}
	if n > max {
		return max
	}
	return n
}

func parseIfMatch(c echo.Context) string {
	raw := strings.TrimSpace(c.Request().Header.Get("If-Match"))
	if raw == "" {
		return ""
	}
	raw = strings.TrimPrefix(raw, "W/")
	raw = strings.TrimSpace(raw)
	if len(raw) >= 2 && raw[0] == '"' {
		raw = strings.Trim(raw, `"`)
	}
	return raw
}

func deletedPtr(ns sql.NullString) *string {
	if !ns.Valid {
		return nil
	}
	s := ns.String
	return &s
}

// RegisterAPIV1 mounts JSON handlers on an Echo group that already uses JWT (+ optional idempotency) middleware.
func (h *Handler) RegisterAPIV1(g *echo.Group) {
	g.GET("/plans", h.APIListPlans)
	g.POST("/plans", h.APICreatePlan)
	g.GET("/plans/:planID", h.APIGetPlan)
	g.PATCH("/plans/:planID", h.APIPatchPlan)
	g.DELETE("/plans/:planID", h.APIDeletePlan)

	g.GET("/plans/:planID/tasks", h.APIListTasks)
	g.GET("/plans/:planID/tasks/week", h.APIListTasksWeek)
	g.POST("/plans/:planID/tasks", h.APICreateTask)
	g.GET("/plans/:planID/access", h.APIListPlanAccess)

	g.GET("/tasks/:taskID", h.APIGetTask)
	g.PATCH("/tasks/:taskID", h.APIPatchTask)
	g.DELETE("/tasks/:taskID", h.APIDeleteTask)
	g.PATCH("/tasks/:taskID/complete", h.APIPatchTaskComplete)

	g.GET("/plans/:planID/resources", h.APIListResources)
	g.POST("/plans/:planID/resources", h.APICreateResource)
	g.GET("/resources/:resourceID", h.APIGetResource)
	g.PATCH("/resources/:resourceID", h.APIPatchResource)
	g.DELETE("/resources/:resourceID", h.APIDeleteResource)

	g.GET("/plans/:planID/templates", h.APIListTemplates)
	g.POST("/plans/:planID/templates", h.APICreateTemplate)
	g.GET("/templates/:templateID", h.APIGetTemplate)
	g.PATCH("/templates/:templateID", h.APIPatchTemplate)
	g.DELETE("/templates/:templateID", h.APIDeleteTemplate)
}

type planDTO struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	User      string  `json:"user"`
	UpdatedAt string  `json:"updated_at"`
	DeletedAt *string `json:"deleted_at,omitempty"`
}

func planToDTO(p database.Plan) planDTO {
	return planDTO{
		ID:        p.ID,
		Name:      p.Name,
		User:      p.User,
		UpdatedAt: p.UpdatedAt,
		DeletedAt: deletedPtr(p.DeletedAt),
	}
}

func (h *Handler) APIListPlans(c echo.Context) error {
	ctx := c.Request().Context()
	uid, err := apiUserID(c)
	if err != nil {
		return apijson.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Missing user context.")
	}

	updatedSince := c.QueryParam("updated_since")
	cursorStr := c.QueryParam("next_cursor")
	limit := int64(parseLimit(c.QueryParam("limit"), 100, 500))

	cursorTs := ""
	var cursorID int64
	if cursorStr != "" {
		k, err := synccursor.Decode(cursorStr)
		if err != nil {
			return apijson.Error(c, http.StatusBadRequest, "BAD_CURSOR", "Invalid next_cursor.")
		}
		cursorTs = k.UpdatedAt
		cursorID = k.ID
	}

	dq := database.New(h.DB)
	rows, err := dq.ListPlansSync(ctx, database.ListPlansSyncParams{
		UserID:       uid,
		UpdatedSince: updatedSince,
		CursorTs:     cursorTs,
		CursorID:     cursorID,
		LimitCount:   limit + 1,
	})
	if err != nil {
		return apijson.ServerError(c, "Could not list plans.", err)
	}

	next := ""
	if len(rows) > int(limit) {
		rows = rows[:limit]
		last := rows[len(rows)-1]
		next = synccursor.Encode(synccursor.Keyset{UpdatedAt: last.UpdatedAt, ID: last.ID})
	}

	out := make([]planDTO, 0, len(rows))
	for _, p := range rows {
		out = append(out, planToDTO(p))
	}
	return c.JSON(http.StatusOK, map[string]any{"items": out, "next_cursor": next})
}

type createPlanBody struct {
	Name string `json:"name"`
}

func (h *Handler) APICreatePlan(c echo.Context) error {
	var body createPlanBody
	if err := c.Bind(&body); err != nil || strings.TrimSpace(body.Name) == "" {
		return apijson.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Expected JSON body with non-empty name.")
	}
	uid, err := apiUserID(c)
	if err != nil {
		return apijson.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Missing user context.")
	}
	ctx := c.Request().Context()
	dq := database.New(h.DB)
	p, err := dq.CreatePlan(ctx, database.CreatePlanParams{Name: body.Name, User: uid})
	if err != nil {
		return apijson.ServerError(c, "Could not create plan.", err)
	}
	return c.JSON(http.StatusCreated, planToDTO(p))
}

func (h *Handler) APIGetPlan(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("planID"), 10, 64)
	if err != nil {
		return apijson.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid plan id.")
	}
	uid, err := apiUserID(c)
	if err != nil {
		return apijson.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Missing user context.")
	}
	ctx := c.Request().Context()
	dq := database.New(h.DB)
	p, err := dq.GetPlan(ctx, database.GetPlanParams{ID: id, UserId: uid})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apijson.Error(c, http.StatusNotFound, "NOT_FOUND", "Plan not found.")
		}
		return apijson.ServerError(c, "Could not load plan.", err)
	}
	return c.JSON(http.StatusOK, planToDTO(p))
}

type patchPlanBody struct {
	Name          string `json:"name"`
	BaseUpdatedAt string `json:"base_updated_at"`
}

func (h *Handler) APIPatchPlan(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("planID"), 10, 64)
	if err != nil {
		return apijson.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid plan id.")
	}
	var body patchPlanBody
	if err := c.Bind(&body); err != nil || strings.TrimSpace(body.Name) == "" {
		return apijson.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Expected JSON with name.")
	}
	uid, err := apiUserID(c)
	if err != nil {
		return apijson.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Missing user context.")
	}
	ctx := c.Request().Context()
	dq := database.New(h.DB)

	base := parseIfMatch(c)
	if base == "" {
		base = body.BaseUpdatedAt
	}

	if base != "" {
		n, err := dq.UpdatePlanIfMatch(ctx, database.UpdatePlanIfMatchParams{
			Name:          body.Name,
			ID:            id,
			UserId:        uid,
			BaseUpdatedAt: base,
		})
		if err != nil {
			return apijson.ServerError(c, "Could not update plan.", err)
		}
		if n == 0 {
			cur, gerr := dq.GetPlan(ctx, database.GetPlanParams{ID: id, UserId: uid})
			if gerr != nil {
				return apijson.Error(c, http.StatusNotFound, "NOT_FOUND", "Plan not found.")
			}
			return apijson.Conflict(c, planToDTO(cur))
		}
	} else {
		if err := dq.UpdatePlan(ctx, database.UpdatePlanParams{Name: body.Name, ID: id, UserId: uid}); err != nil {
			return apijson.ServerError(c, "Could not update plan.", err)
		}
	}
	p, err := dq.GetPlan(ctx, database.GetPlanParams{ID: id, UserId: uid})
	if err != nil {
		return apijson.ServerError(c, "Could not load plan.", err)
	}
	return c.JSON(http.StatusOK, planToDTO(p))
}

func (h *Handler) APIDeletePlan(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("planID"), 10, 64)
	if err != nil {
		return apijson.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid plan id.")
	}
	uid, err := apiUserID(c)
	if err != nil {
		return apijson.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Missing user context.")
	}
	ctx := c.Request().Context()
	dq := database.New(h.DB)
	if err := dq.DeletePlan(ctx, database.DeletePlanParams{ID: id, UserId: uid}); err != nil {
		return apijson.ServerError(c, "Could not delete plan.", err)
	}
	return c.NoContent(http.StatusNoContent)
}
