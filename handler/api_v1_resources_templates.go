package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/sknutsen/planner/database"
	"github.com/sknutsen/planner/internal/apijson"
	"github.com/sknutsen/planner/internal/synccursor"
	"github.com/sknutsen/planner/lib"
)

type resourceDTO struct {
	ID           int64           `json:"id"`
	PlanID       int64           `json:"plan_id"`
	Title        string          `json:"title"`
	ResourceType int64           `json:"resource_type"`
	Content      json.RawMessage `json:"content"`
	UpdatedAt    string          `json:"updated_at"`
	DeletedAt    *string         `json:"deleted_at,omitempty"`
}

func contentJSON(c interface{}) json.RawMessage {
	if c == nil {
		return json.RawMessage(`null`)
	}
	switch v := c.(type) {
	case []byte:
		if len(v) == 0 {
			return json.RawMessage(`null`)
		}
		return json.RawMessage(v)
	case string:
		if v == "" {
			return json.RawMessage(`null`)
		}
		return json.RawMessage(v)
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return json.RawMessage(`null`)
		}
		return json.RawMessage(b)
	}
}

func resourceToDTO(r database.Resource) resourceDTO {
	return resourceDTO{
		ID:           r.ID,
		PlanID:       r.PlanID,
		Title:        r.Title,
		ResourceType: r.ResourceType,
		Content:      contentJSON(r.Content),
		UpdatedAt:    r.UpdatedAt,
		DeletedAt:    deletedPtr(r.DeletedAt),
	}
}

func (h *Handler) APIListResources(c echo.Context) error {
	planID, err := strconv.ParseInt(c.Param("planID"), 10, 64)
	if err != nil {
		return apijson.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid plan id.")
	}
	uid, err := apiUserID(c)
	if err != nil {
		return apijson.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Missing user context.")
	}
	ctx := c.Request().Context()

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
	rows, err := dq.ListResourcesByPlanSync(ctx, database.ListResourcesByPlanSyncParams{
		PlanID:       planID,
		UserID:       uid,
		UpdatedSince: updatedSince,
		CursorTs:     cursorTs,
		CursorID:     cursorID,
		LimitCount:   limit + 1,
	})
	if err != nil {
		return apijson.ServerError(c, "Could not list resources.", err)
	}

	next := ""
	if len(rows) > int(limit) {
		rows = rows[:limit]
		last := rows[len(rows)-1]
		next = synccursor.Encode(synccursor.Keyset{UpdatedAt: last.UpdatedAt, ID: last.ID})
	}
	out := make([]resourceDTO, 0, len(rows))
	for _, r := range rows {
		out = append(out, resourceToDTO(r))
	}
	return c.JSON(http.StatusOK, map[string]any{"items": out, "next_cursor": next})
}

type createResourceBody struct {
	Title        string          `json:"title"`
	ResourceType int64           `json:"resource_type"`
	Content      json.RawMessage `json:"content"`
}

func (h *Handler) APICreateResource(c echo.Context) error {
	planID, err := strconv.ParseInt(c.Param("planID"), 10, 64)
	if err != nil {
		return apijson.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid plan id.")
	}
	var body createResourceBody
	if err := c.Bind(&body); err != nil || strings.TrimSpace(body.Title) == "" {
		return apijson.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Expected title.")
	}
	var content interface{}
	if len(body.Content) > 0 {
		content = string(body.Content)
	}
	if _, err := apiUserID(c); err != nil {
		return apijson.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Missing user context.")
	}
	ctx := c.Request().Context()
	dq := database.New(h.DB)
	r, err := dq.CreateResource(ctx, database.CreateResourceParams{
		PlanID:       planID,
		Title:        body.Title,
		ResourceType: body.ResourceType,
		Content:      content,
	})
	if err != nil {
		return apijson.ServerError(c, "Could not create resource.", err)
	}
	return c.JSON(http.StatusCreated, resourceToDTO(r))
}

func (h *Handler) APIGetResource(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("resourceID"), 10, 64)
	if err != nil {
		return apijson.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid resource id.")
	}
	uid, err := apiUserID(c)
	if err != nil {
		return apijson.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Missing user context.")
	}
	ctx := c.Request().Context()
	dq := database.New(h.DB)
	r, err := dq.GetResource(ctx, database.GetResourceParams{ID: id, UserId: uid})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apijson.Error(c, http.StatusNotFound, "NOT_FOUND", "Resource not found.")
		}
		return apijson.ServerError(c, "Could not load resource.", err)
	}
	return c.JSON(http.StatusOK, resourceToDTO(r))
}

type patchResourceBody struct {
	Title         string          `json:"title"`
	ResourceType  int64           `json:"resource_type"`
	Content       json.RawMessage `json:"content"`
	BaseUpdatedAt string          `json:"base_updated_at"`
}

func (h *Handler) APIPatchResource(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("resourceID"), 10, 64)
	if err != nil {
		return apijson.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid resource id.")
	}
	var body patchResourceBody
	if err := c.Bind(&body); err != nil || strings.TrimSpace(body.Title) == "" {
		return apijson.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Expected title.")
	}
	var content interface{}
	if len(body.Content) > 0 {
		content = string(body.Content)
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
		n, err := dq.UpdateResourceIfMatch(ctx, database.UpdateResourceIfMatchParams{
			Title:         body.Title,
			ResourceType:  body.ResourceType,
			Content:       content,
			ID:            id,
			UserId:        uid,
			BaseUpdatedAt: base,
		})
		if err != nil {
			return apijson.ServerError(c, "Could not update resource.", err)
		}
		if n == 0 {
			cur, gerr := dq.GetResource(ctx, database.GetResourceParams{ID: id, UserId: uid})
			if gerr != nil {
				return apijson.Error(c, http.StatusNotFound, "NOT_FOUND", "Resource not found.")
			}
			return apijson.Conflict(c, resourceToDTO(cur))
		}
	} else {
		if err := dq.UpdateResource(ctx, database.UpdateResourceParams{
			Title:        body.Title,
			ResourceType: body.ResourceType,
			Content:      content,
			ID:           id,
			UserId:       uid,
		}); err != nil {
			return apijson.ServerError(c, "Could not update resource.", err)
		}
	}
	r, err := dq.GetResource(ctx, database.GetResourceParams{ID: id, UserId: uid})
	if err != nil {
		return apijson.ServerError(c, "Could not load resource.", err)
	}
	return c.JSON(http.StatusOK, resourceToDTO(r))
}

func (h *Handler) APIDeleteResource(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("resourceID"), 10, 64)
	if err != nil {
		return apijson.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid resource id.")
	}
	uid, err := apiUserID(c)
	if err != nil {
		return apijson.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Missing user context.")
	}
	ctx := c.Request().Context()
	dq := database.New(h.DB)
	if err := dq.DeleteResource(ctx, database.DeleteResourceParams{ID: id, UserId: uid}); err != nil {
		return apijson.ServerError(c, "Could not delete resource.", err)
	}
	return c.NoContent(http.StatusNoContent)
}

type templateDTO struct {
	ID          int64   `json:"id"`
	PlanID      int64   `json:"plan_id"`
	Title       string  `json:"title"`
	Subtitle    string  `json:"subtitle,omitempty"`
	Description string  `json:"description,omitempty"`
	UpdatedAt   string  `json:"updated_at"`
	DeletedAt   *string `json:"deleted_at,omitempty"`
}

func templateToDTO(t database.Template) templateDTO {
	return templateDTO{
		ID:          t.ID,
		PlanID:      t.PlanID,
		Title:       t.Title,
		Subtitle:    lib.AsString(t.Subtitle),
		Description: lib.AsString(t.Description),
		UpdatedAt:   t.UpdatedAt,
		DeletedAt:   deletedPtr(t.DeletedAt),
	}
}

func (h *Handler) APIListTemplates(c echo.Context) error {
	planID, err := strconv.ParseInt(c.Param("planID"), 10, 64)
	if err != nil {
		return apijson.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid plan id.")
	}
	uid, err := apiUserID(c)
	if err != nil {
		return apijson.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Missing user context.")
	}
	ctx := c.Request().Context()

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
	rows, err := dq.ListTemplatesByPlanSync(ctx, database.ListTemplatesByPlanSyncParams{
		PlanID:       planID,
		UserID:       uid,
		UpdatedSince: updatedSince,
		CursorTs:     cursorTs,
		CursorID:     cursorID,
		LimitCount:   limit + 1,
	})
	if err != nil {
		return apijson.ServerError(c, "Could not list templates.", err)
	}
	next := ""
	if len(rows) > int(limit) {
		rows = rows[:limit]
		last := rows[len(rows)-1]
		next = synccursor.Encode(synccursor.Keyset{UpdatedAt: last.UpdatedAt, ID: last.ID})
	}
	out := make([]templateDTO, 0, len(rows))
	for _, t := range rows {
		out = append(out, templateToDTO(t))
	}
	return c.JSON(http.StatusOK, map[string]any{"items": out, "next_cursor": next})
}

type createTemplateBody struct {
	Title       string `json:"title"`
	Subtitle    string `json:"subtitle"`
	Description string `json:"description"`
}

func (h *Handler) APICreateTemplate(c echo.Context) error {
	planID, err := strconv.ParseInt(c.Param("planID"), 10, 64)
	if err != nil {
		return apijson.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid plan id.")
	}
	var body createTemplateBody
	if err := c.Bind(&body); err != nil || strings.TrimSpace(body.Title) == "" {
		return apijson.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Expected title.")
	}
	if _, err := apiUserID(c); err != nil {
		return apijson.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Missing user context.")
	}
	ctx := c.Request().Context()
	dq := database.New(h.DB)
	t, err := dq.CreateTemplate(ctx, database.CreateTemplateParams{
		PlanID:      planID,
		Title:       body.Title,
		Subtitle:    body.Subtitle,
		Description: body.Description,
	})
	if err != nil {
		return apijson.ServerError(c, "Could not create template.", err)
	}
	return c.JSON(http.StatusCreated, templateToDTO(t))
}

func (h *Handler) APIGetTemplate(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("templateID"), 10, 64)
	if err != nil {
		return apijson.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid template id.")
	}
	uid, err := apiUserID(c)
	if err != nil {
		return apijson.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Missing user context.")
	}
	ctx := c.Request().Context()
	dq := database.New(h.DB)
	t, err := dq.GetTemplate(ctx, database.GetTemplateParams{ID: id, UserId: uid})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apijson.Error(c, http.StatusNotFound, "NOT_FOUND", "Template not found.")
		}
		return apijson.ServerError(c, "Could not load template.", err)
	}
	return c.JSON(http.StatusOK, templateToDTO(t))
}

type patchTemplateBody struct {
	Title         string `json:"title"`
	Subtitle      string `json:"subtitle"`
	Description   string `json:"description"`
	BaseUpdatedAt string `json:"base_updated_at"`
}

func (h *Handler) APIPatchTemplate(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("templateID"), 10, 64)
	if err != nil {
		return apijson.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid template id.")
	}
	var body patchTemplateBody
	if err := c.Bind(&body); err != nil || strings.TrimSpace(body.Title) == "" {
		return apijson.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Expected title.")
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
		n, err := dq.UpdateTemplateIfMatch(ctx, database.UpdateTemplateIfMatchParams{
			Title:         body.Title,
			Subtitle:      body.Subtitle,
			Description:   body.Description,
			ID:            id,
			UserId:        uid,
			BaseUpdatedAt: base,
		})
		if err != nil {
			return apijson.ServerError(c, "Could not update template.", err)
		}
		if n == 0 {
			cur, gerr := dq.GetTemplate(ctx, database.GetTemplateParams{ID: id, UserId: uid})
			if gerr != nil {
				return apijson.Error(c, http.StatusNotFound, "NOT_FOUND", "Template not found.")
			}
			return apijson.Conflict(c, templateToDTO(cur))
		}
	} else {
		if err := dq.UpdateTemplate(ctx, database.UpdateTemplateParams{
			Title:       body.Title,
			Subtitle:    body.Subtitle,
			Description: body.Description,
			ID:          id,
			UserId:      uid,
		}); err != nil {
			return apijson.ServerError(c, "Could not update template.", err)
		}
	}
	t, err := dq.GetTemplate(ctx, database.GetTemplateParams{ID: id, UserId: uid})
	if err != nil {
		return apijson.ServerError(c, "Could not load template.", err)
	}
	return c.JSON(http.StatusOK, templateToDTO(t))
}

func (h *Handler) APIDeleteTemplate(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("templateID"), 10, 64)
	if err != nil {
		return apijson.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid template id.")
	}
	uid, err := apiUserID(c)
	if err != nil {
		return apijson.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Missing user context.")
	}
	ctx := c.Request().Context()
	dq := database.New(h.DB)
	if err := dq.DeleteTemplate(ctx, database.DeleteTemplateParams{ID: id, UserId: uid}); err != nil {
		return apijson.ServerError(c, "Could not delete template.", err)
	}
	return c.NoContent(http.StatusNoContent)
}

type planAccessDTO struct {
	ID        int64   `json:"id"`
	PlanID    int64   `json:"plan_id"`
	User      string  `json:"user"`
	UpdatedAt string  `json:"updated_at"`
	DeletedAt *string `json:"deleted_at,omitempty"`
}

func planAccessToDTO(pa database.PlanAccess) planAccessDTO {
	return planAccessDTO{
		ID:        pa.ID,
		PlanID:    pa.PlanID,
		User:      pa.User,
		UpdatedAt: pa.UpdatedAt,
		DeletedAt: deletedPtr(pa.DeletedAt),
	}
}

func (h *Handler) APIListPlanAccess(c echo.Context) error {
	planID, err := strconv.ParseInt(c.Param("planID"), 10, 64)
	if err != nil {
		return apijson.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid plan id.")
	}
	uid, err := apiUserID(c)
	if err != nil {
		return apijson.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Missing user context.")
	}
	ctx := c.Request().Context()

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
	rows, err := dq.ListPlanAccessByPlanSync(ctx, database.ListPlanAccessByPlanSyncParams{
		PlanID:       planID,
		UserID:       uid,
		UpdatedSince: updatedSince,
		CursorTs:     cursorTs,
		CursorID:     cursorID,
		LimitCount:   limit + 1,
	})
	if err != nil {
		return apijson.ServerError(c, "Could not list plan access rows.", err)
	}
	next := ""
	if len(rows) > int(limit) {
		rows = rows[:limit]
		last := rows[len(rows)-1]
		next = synccursor.Encode(synccursor.Keyset{UpdatedAt: last.UpdatedAt, ID: last.ID})
	}
	out := make([]planAccessDTO, 0, len(rows))
	for _, pa := range rows {
		out = append(out, planAccessToDTO(pa))
	}
	return c.JSON(http.StatusOK, map[string]any{"items": out, "next_cursor": next})
}
