package handler

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
	"github.com/sknutsen/planner/lib"
)

type taskDTO struct {
	ID          int64   `json:"id"`
	PlanID      int64   `json:"plan_id"`
	Date        string  `json:"date"`
	Title       string  `json:"title"`
	Subtitle    string  `json:"subtitle,omitempty"`
	Description string  `json:"description,omitempty"`
	IsComplete  bool    `json:"is_complete"`
	UpdatedAt   string  `json:"updated_at"`
	DeletedAt   *string `json:"deleted_at,omitempty"`
}

func taskToDTO(t database.Task) taskDTO {
	return taskDTO{
		ID:          t.ID,
		PlanID:      t.PlanID,
		Date:        lib.StripDateString(t.Date),
		Title:       t.Title,
		Subtitle:    lib.AsString(t.Subtitle),
		Description: lib.AsString(t.Description),
		IsComplete:  t.IsComplete != 0,
		UpdatedAt:   t.UpdatedAt,
		DeletedAt:   deletedPtr(t.DeletedAt),
	}
}

func (h *Handler) APIListTasks(c echo.Context) error {
	planID, err := strconv.ParseInt(c.Param("planID"), 10, 64)
	if err != nil {
		return apijson.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid plan id.")
	}
	uid, err := apiUserID(c)
	if err != nil {
		return apijson.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Missing user context.")
	}
	ctx := c.Request().Context()

	date := c.QueryParam("date")
	if date != "" {
		dq := database.New(h.DB)
		tasks, err := dq.GetTasksByDate(ctx, database.GetTasksByDateParams{
			Date:   date,
			PlanId: planID,
			UserId: uid,
		})
		if err != nil {
			return apijson.Error(c, http.StatusInternalServerError, "SERVER_ERROR", "Could not list tasks.")
		}
		out := make([]taskDTO, 0, len(tasks))
		for _, t := range tasks {
			out = append(out, taskToDTO(t))
		}
		return c.JSON(http.StatusOK, map[string]any{"items": out})
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
	rows, err := dq.ListTasksByPlanSync(ctx, database.ListTasksByPlanSyncParams{
		PlanID:       planID,
		UserID:       uid,
		UpdatedSince: updatedSince,
		CursorTs:     cursorTs,
		CursorID:     cursorID,
		LimitCount:   limit + 1,
	})
	if err != nil {
		return apijson.Error(c, http.StatusInternalServerError, "SERVER_ERROR", "Could not list tasks.")
	}
	next := ""
	if len(rows) > int(limit) {
		rows = rows[:limit]
		last := rows[len(rows)-1]
		next = synccursor.Encode(synccursor.Keyset{UpdatedAt: last.UpdatedAt, ID: last.ID})
	}
	out := make([]taskDTO, 0, len(rows))
	for _, t := range rows {
		out = append(out, taskToDTO(t))
	}
	return c.JSON(http.StatusOK, map[string]any{"items": out, "next_cursor": next})
}

func (h *Handler) APIListTasksWeek(c echo.Context) error {
	planID, err := strconv.ParseInt(c.Param("planID"), 10, 64)
	if err != nil {
		return apijson.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid plan id.")
	}
	week := c.QueryParam("week")
	if week == "" {
		return apijson.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Query week is required (e.g. 2026-20).")
	}
	y, w := lib.ISOWeekFromString(week)
	dates := lib.DatesInWeek(y, w)
	uid, err := apiUserID(c)
	if err != nil {
		return apijson.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Missing user context.")
	}
	ctx := c.Request().Context()
	dq := database.New(h.DB)
	tasks, err := dq.ListTasksByPlanAndDates(ctx, database.ListTasksByPlanAndDatesParams{
		PlanID: planID,
		UserID: uid,
		Dates:  dates,
	})
	if err != nil {
		return apijson.Error(c, http.StatusInternalServerError, "SERVER_ERROR", "Could not list week tasks.")
	}
	out := make([]taskDTO, 0, len(tasks))
	for _, t := range tasks {
		out = append(out, taskToDTO(t))
	}
	return c.JSON(http.StatusOK, map[string]any{"week": week, "items": out})
}

type createTaskBody struct {
	Title       string `json:"title"`
	Date        string `json:"date"`
	Subtitle    string `json:"subtitle"`
	Description string `json:"description"`
	TemplateID  int64  `json:"template_id"`
}

func (h *Handler) APICreateTask(c echo.Context) error {
	planID, err := strconv.ParseInt(c.Param("planID"), 10, 64)
	if err != nil {
		return apijson.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid plan id.")
	}
	var body createTaskBody
	if err := c.Bind(&body); err != nil || strings.TrimSpace(body.Title) == "" || strings.TrimSpace(body.Date) == "" {
		return apijson.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Expected title and date.")
	}
	uid, err := apiUserID(c)
	if err != nil {
		return apijson.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Missing user context.")
	}
	ctx := c.Request().Context()
	dq := database.New(h.DB)
	if body.TemplateID != 0 {
		t, err := dq.CreateTaskFromTemplate(ctx, database.CreateTaskFromTemplateParams{
			Date:       body.Date,
			TemplateId: body.TemplateID,
			UserId:     uid,
		})
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return apijson.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Template not found or inaccessible.")
			}
			return apijson.Error(c, http.StatusInternalServerError, "SERVER_ERROR", "Could not create task from template.")
		}
		return c.JSON(http.StatusCreated, taskToDTO(t))
	}
	t, err := dq.CreateTask(ctx, database.CreateTaskParams{
		PlanID:      planID,
		Title:       body.Title,
		Date:        body.Date,
		Subtitle:    body.Subtitle,
		Description: body.Description,
	})
	if err != nil {
		return apijson.Error(c, http.StatusInternalServerError, "SERVER_ERROR", "Could not create task.")
	}
	return c.JSON(http.StatusCreated, taskToDTO(t))
}

func (h *Handler) APIGetTask(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("taskID"), 10, 64)
	if err != nil {
		return apijson.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid task id.")
	}
	uid, err := apiUserID(c)
	if err != nil {
		return apijson.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Missing user context.")
	}
	ctx := c.Request().Context()
	dq := database.New(h.DB)
	t, err := dq.GetTask(ctx, database.GetTaskParams{ID: id, UserId: uid})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apijson.Error(c, http.StatusNotFound, "NOT_FOUND", "Task not found.")
		}
		return apijson.Error(c, http.StatusInternalServerError, "SERVER_ERROR", "Could not load task.")
	}
	return c.JSON(http.StatusOK, taskToDTO(t))
}

type patchTaskBody struct {
	Title         string `json:"title"`
	Date          string `json:"date"`
	Subtitle      string `json:"subtitle"`
	Description   string `json:"description"`
	BaseUpdatedAt string `json:"base_updated_at"`
}

func (h *Handler) APIPatchTask(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("taskID"), 10, 64)
	if err != nil {
		return apijson.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid task id.")
	}
	var body patchTaskBody
	if err := c.Bind(&body); err != nil {
		return apijson.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON body.")
	}
	if strings.TrimSpace(body.Title) == "" || strings.TrimSpace(body.Date) == "" {
		return apijson.Error(c, http.StatusBadRequest, "BAD_REQUEST", "title and date are required.")
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
		n, err := dq.UpdateTaskIfMatch(ctx, database.UpdateTaskIfMatchParams{
			Title:         body.Title,
			Subtitle:      body.Subtitle,
			Date:          body.Date,
			Description:   body.Description,
			ID:            id,
			UserId:        uid,
			BaseUpdatedAt: base,
		})
		if err != nil {
			return apijson.Error(c, http.StatusInternalServerError, "SERVER_ERROR", "Could not update task.")
		}
		if n == 0 {
			cur, gerr := dq.GetTask(ctx, database.GetTaskParams{ID: id, UserId: uid})
			if gerr != nil {
				return apijson.Error(c, http.StatusNotFound, "NOT_FOUND", "Task not found.")
			}
			return apijson.Conflict(c, taskToDTO(cur))
		}
	} else {
		err := dq.UpdateTask(ctx, database.UpdateTaskParams{
			Title:       body.Title,
			Subtitle:    body.Subtitle,
			Date:        body.Date,
			Description: body.Description,
			ID:          id,
			UserId:      uid,
		})
		if err != nil {
			return apijson.Error(c, http.StatusInternalServerError, "SERVER_ERROR", "Could not update task.")
		}
	}
	t, err := dq.GetTask(ctx, database.GetTaskParams{ID: id, UserId: uid})
	if err != nil {
		return apijson.Error(c, http.StatusInternalServerError, "SERVER_ERROR", "Could not load task.")
	}
	return c.JSON(http.StatusOK, taskToDTO(t))
}

func (h *Handler) APIDeleteTask(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("taskID"), 10, 64)
	if err != nil {
		return apijson.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid task id.")
	}
	uid, err := apiUserID(c)
	if err != nil {
		return apijson.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Missing user context.")
	}
	ctx := c.Request().Context()
	dq := database.New(h.DB)
	if err := dq.DeleteTask(ctx, database.DeleteTaskParams{ID: id, UserId: uid}); err != nil {
		return apijson.Error(c, http.StatusInternalServerError, "SERVER_ERROR", "Could not delete task.")
	}
	return c.NoContent(http.StatusNoContent)
}

type patchCompleteBody struct {
	IsComplete    bool   `json:"is_complete"`
	BaseUpdatedAt string `json:"base_updated_at"`
}

func (h *Handler) APIPatchTaskComplete(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("taskID"), 10, 64)
	if err != nil {
		return apijson.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid task id.")
	}
	var body patchCompleteBody
	if err := c.Bind(&body); err != nil {
		return apijson.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON body.")
	}
	uid, err := apiUserID(c)
	if err != nil {
		return apijson.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Missing user context.")
	}
	ctx := c.Request().Context()
	dq := database.New(h.DB)

	var ic int64
	if body.IsComplete {
		ic = 1
	}
	base := parseIfMatch(c)
	if base == "" {
		base = body.BaseUpdatedAt
	}
	if base != "" {
		cur, err := dq.GetTask(ctx, database.GetTaskParams{ID: id, UserId: uid})
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return apijson.Error(c, http.StatusNotFound, "NOT_FOUND", "Task not found.")
			}
			return apijson.Error(c, http.StatusInternalServerError, "SERVER_ERROR", "Could not load task.")
		}
		if cur.UpdatedAt != base {
			return apijson.Conflict(c, taskToDTO(cur))
		}
	}
	if err := dq.SetIsCompleteTask(ctx, database.SetIsCompleteTaskParams{
		IsComplete: ic,
		ID:         id,
		UserId:     uid,
	}); err != nil {
		return apijson.Error(c, http.StatusInternalServerError, "SERVER_ERROR", "Could not update task.")
	}
	t, err := dq.GetTask(ctx, database.GetTaskParams{ID: id, UserId: uid})
	if err != nil {
		return apijson.Error(c, http.StatusInternalServerError, "SERVER_ERROR", "Could not load task.")
	}
	return c.JSON(http.StatusOK, taskToDTO(t))
}
