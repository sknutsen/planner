package database

import (
	"context"
	"database/sql"
	"errors"
	"testing"
)

const (
	ownerSub    = "auth0|owner"
	guestSub    = "auth0|guest"
	outsiderSub = "auth0|outsider"
)

type seedData struct {
	planID     int64
	taskID     int64
	resourceID int64
	templateID int64
	accessID   int64
}

func seedPlanner(t *testing.T, q *Queries) seedData {
	t.Helper()
	ctx := context.Background()

	plan, err := q.CreatePlan(ctx, CreatePlanParams{Name: "Week plan", User: ownerSub})
	if err != nil {
		t.Fatalf("CreatePlan: %v", err)
	}
	if err := q.GrantAccess(ctx, GrantAccessParams{PlanID: plan.ID, User: guestSub}); err != nil {
		t.Fatalf("GrantAccess: %v", err)
	}
	accessRows, err := q.ListPlanAccess(ctx, ListPlanAccessParams{ID: plan.ID, User: ownerSub})
	if err != nil || len(accessRows) != 1 {
		t.Fatalf("ListPlanAccess: rows=%d err=%v", len(accessRows), err)
	}

	task, err := q.CreateTask(ctx, CreateTaskParams{
		PlanID:      plan.ID,
		Title:       "Task A",
		Date:        "2026-06-13",
		Subtitle:    "sub",
		Description: "desc",
	})
	if err != nil {
		t.Fatalf("CreateTask: %v", err)
	}

	resource, err := q.CreateResource(ctx, CreateResourceParams{
		PlanID:       plan.ID,
		Title:        "Note",
		ResourceType: 1,
		Content:      `{"body":"hi"}`,
	})
	if err != nil {
		t.Fatalf("CreateResource: %v", err)
	}

	tmpl, err := q.CreateTemplate(ctx, CreateTemplateParams{
		PlanID:      plan.ID,
		Title:       "Template A",
		Subtitle:    "t-sub",
		Description: "t-desc",
	})
	if err != nil {
		t.Fatalf("CreateTemplate: %v", err)
	}

	return seedData{
		planID:     plan.ID,
		taskID:     task.ID,
		resourceID: resource.ID,
		templateID: tmpl.ID,
		accessID:   accessRows[0].ID,
	}
}

func TestSchemaMigrationsApply(t *testing.T) {
	openTestDB(t)
}

func TestPlansQueries(t *testing.T) {
	db := openTestDB(t)
	q := New(db)
	ctx := context.Background()
	s := seedPlanner(t, q)

	got, err := q.GetPlan(ctx, GetPlanParams{ID: s.planID, UserId: ownerSub})
	if err != nil || got.ID != s.planID {
		t.Fatalf("GetPlan owner: %+v err=%v", got, err)
	}
	if _, err := q.GetPlan(ctx, GetPlanParams{ID: s.planID, UserId: outsiderSub}); !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("GetPlan outsider want ErrNoRows, got %v", err)
	}
	if _, err := q.GetPlan(ctx, GetPlanParams{ID: s.planID, UserId: guestSub}); err != nil {
		t.Fatalf("GetPlan guest: %v", err)
	}

	plans, err := q.ListPlans(ctx, ownerSub)
	if err != nil || len(plans) != 1 {
		t.Fatalf("ListPlans: %d err=%v", len(plans), err)
	}

	syncRows, err := q.ListPlansSync(ctx, ListPlansSyncParams{
		UserID:       ownerSub,
		UpdatedSince: "",
		CursorTs:     "",
		LimitCount:   10,
	})
	if err != nil || len(syncRows) != 1 {
		t.Fatalf("ListPlansSync: %d err=%v", len(syncRows), err)
	}
	page2, err := q.ListPlansSync(ctx, ListPlansSyncParams{
		UserID:       ownerSub,
		UpdatedSince: "",
		CursorTs:     syncRows[0].UpdatedAt,
		CursorID:     syncRows[0].ID,
		LimitCount:   10,
	})
	if err != nil {
		t.Fatalf("ListPlansSync cursor: %v", err)
	}
	if len(page2) != 0 {
		t.Fatalf("ListPlansSync cursor: want empty page, got %d", len(page2))
	}

	if err := q.UpdatePlan(ctx, UpdatePlanParams{Name: "Renamed", ID: s.planID, UserId: ownerSub}); err != nil {
		t.Fatalf("UpdatePlan: %v", err)
	}
	updated, err := q.GetPlan(ctx, GetPlanParams{ID: s.planID, UserId: ownerSub})
	if err != nil || updated.Name != "Renamed" {
		t.Fatalf("UpdatePlan result: %+v err=%v", updated, err)
	}

	stale := updated.UpdatedAt
	n, err := q.UpdatePlanIfMatch(ctx, UpdatePlanIfMatchParams{
		Name:          "Nope",
		ID:            s.planID,
		UserId:        ownerSub,
		BaseUpdatedAt: "1970-01-01T00:00:00Z",
	})
	if err != nil || n != 0 {
		t.Fatalf("UpdatePlanIfMatch stale: n=%d err=%v", n, err)
	}
	n, err = q.UpdatePlanIfMatch(ctx, UpdatePlanIfMatchParams{
		Name:          "Matched",
		ID:            s.planID,
		UserId:        ownerSub,
		BaseUpdatedAt: stale,
	})
	if err != nil || n != 1 {
		t.Fatalf("UpdatePlanIfMatch match: n=%d err=%v", n, err)
	}
}

func TestTasksQueries(t *testing.T) {
	db := openTestDB(t)
	q := New(db)
	ctx := context.Background()
	s := seedPlanner(t, q)

	task, err := q.GetTask(ctx, GetTaskParams{ID: s.taskID, UserId: guestSub})
	if err != nil || task.Title != "Task A" {
		t.Fatalf("GetTask guest: %+v err=%v", task, err)
	}
	if _, err := q.GetTask(ctx, GetTaskParams{ID: s.taskID, UserId: outsiderSub}); !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("GetTask outsider: %v", err)
	}

	byDate, err := q.GetTasksByDate(ctx, GetTasksByDateParams{
		Date:   "2026-06-13",
		PlanId: s.planID,
		UserId: ownerSub,
	})
	if err != nil || len(byDate) != 1 {
		t.Fatalf("GetTasksByDate: %d err=%v", len(byDate), err)
	}

	byPlan, err := q.GetTasksByPlan(ctx, GetTasksByPlanParams{PlanId: s.planID, UserId: ownerSub})
	if err != nil || len(byPlan) != 1 {
		t.Fatalf("GetTasksByPlan: %d err=%v", len(byPlan), err)
	}

	syncRows, err := q.ListTasksByPlanSync(ctx, ListTasksByPlanSyncParams{
		PlanID:       s.planID,
		UserID:       ownerSub,
		UpdatedSince: "",
		CursorTs:     "",
		LimitCount:   10,
	})
	if err != nil || len(syncRows) != 1 {
		t.Fatalf("ListTasksByPlanSync: %d err=%v", len(syncRows), err)
	}

	weekRows, err := q.ListTasksByPlanAndDates(ctx, ListTasksByPlanAndDatesParams{
		PlanID: s.planID,
		UserID: ownerSub,
		Dates:  []string{"2026-06-13", "2026-06-14"},
	})
	if err != nil || len(weekRows) != 1 {
		t.Fatalf("ListTasksByPlanAndDates: %d err=%v", len(weekRows), err)
	}

	fromTmpl, err := q.CreateTaskFromTemplate(ctx, CreateTaskFromTemplateParams{
		Date:       "2026-06-14",
		TemplateId: s.templateID,
		UserId:     guestSub,
	})
	if err != nil || fromTmpl.Date != "2026-06-14" {
		t.Fatalf("CreateTaskFromTemplate: %+v err=%v", fromTmpl, err)
	}
	if _, err := q.CreateTaskFromTemplate(ctx, CreateTaskFromTemplateParams{
		Date:       "2026-06-15",
		TemplateId: 99999,
		UserId:     ownerSub,
	}); !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("CreateTaskFromTemplate missing: %v", err)
	}

	if err := q.SetIsCompleteTask(ctx, SetIsCompleteTaskParams{IsComplete: 1, ID: s.taskID, UserId: guestSub}); err != nil {
		t.Fatalf("SetIsCompleteTask: %v", err)
	}
	if err := q.UpdateTask(ctx, UpdateTaskParams{
		Title: "Updated", Subtitle: "s", Date: "2026-06-13", Description: "d",
		ID: s.taskID, UserId: ownerSub,
	}); err != nil {
		t.Fatalf("UpdateTask: %v", err)
	}
	cur, err := q.GetTask(ctx, GetTaskParams{ID: s.taskID, UserId: ownerSub})
	if err != nil {
		t.Fatalf("GetTask after update: %v", err)
	}
	n, err := q.UpdateTaskIfMatch(ctx, UpdateTaskIfMatchParams{
		Title: "X", Subtitle: "s", Date: "2026-06-13", Description: "d",
		ID: s.taskID, UserId: ownerSub, BaseUpdatedAt: cur.UpdatedAt,
	})
	if err != nil || n != 1 {
		t.Fatalf("UpdateTaskIfMatch: n=%d err=%v", n, err)
	}
}

func TestResourcesAndTemplatesQueries(t *testing.T) {
	db := openTestDB(t)
	q := New(db)
	ctx := context.Background()
	s := seedPlanner(t, q)

	res, err := q.GetResource(ctx, GetResourceParams{ID: s.resourceID, UserId: guestSub})
	if err != nil || res.Title != "Note" {
		t.Fatalf("GetResource: %+v err=%v", res, err)
	}
	resources, err := q.GetResourcesByPlan(ctx, GetResourcesByPlanParams{PlanId: s.planID, UserId: ownerSub})
	if err != nil || len(resources) != 1 {
		t.Fatalf("GetResourcesByPlan: %d err=%v", len(resources), err)
	}
	resSync, err := q.ListResourcesByPlanSync(ctx, ListResourcesByPlanSyncParams{
		PlanID: s.planID, UserID: ownerSub, UpdatedSince: "", CursorTs: "", LimitCount: 10,
	})
	if err != nil || len(resSync) != 1 {
		t.Fatalf("ListResourcesByPlanSync: %d err=%v", len(resSync), err)
	}
	if err := q.UpdateResource(ctx, UpdateResourceParams{
		Title: "Note 2", ResourceType: 2, Content: `{}`,
		ID: s.resourceID, UserId: ownerSub,
	}); err != nil {
		t.Fatalf("UpdateResource: %v", err)
	}
	cur, err := q.GetResource(ctx, GetResourceParams{ID: s.resourceID, UserId: ownerSub})
	if err != nil {
		t.Fatalf("GetResource reload: %v", err)
	}
	if n, err := q.UpdateResourceIfMatch(ctx, UpdateResourceIfMatchParams{
		Title: "Note 3", ResourceType: 2, Content: `{}`,
		ID: s.resourceID, UserId: ownerSub, BaseUpdatedAt: cur.UpdatedAt,
	}); err != nil || n != 1 {
		t.Fatalf("UpdateResourceIfMatch: n=%d err=%v", n, err)
	}

	tmpl, err := q.GetTemplate(ctx, GetTemplateParams{ID: s.templateID, UserId: guestSub})
	if err != nil || tmpl.Title != "Template A" {
		t.Fatalf("GetTemplate: %+v err=%v", tmpl, err)
	}
	templates, err := q.GetTemplatesByPlan(ctx, GetTemplatesByPlanParams{PlanId: s.planID, UserId: ownerSub})
	if err != nil || len(templates) != 1 {
		t.Fatalf("GetTemplatesByPlan: %d err=%v", len(templates), err)
	}
	tmplSync, err := q.ListTemplatesByPlanSync(ctx, ListTemplatesByPlanSyncParams{
		PlanID: s.planID, UserID: ownerSub, UpdatedSince: "", CursorTs: "", LimitCount: 10,
	})
	if err != nil || len(tmplSync) != 1 {
		t.Fatalf("ListTemplatesByPlanSync: %d err=%v", len(tmplSync), err)
	}
	if err := q.UpdateTemplate(ctx, UpdateTemplateParams{
		Title: "T2", Subtitle: "s", Description: "d",
		ID: s.templateID, UserId: ownerSub,
	}); err != nil {
		t.Fatalf("UpdateTemplate: %v", err)
	}
	curTmpl, err := q.GetTemplate(ctx, GetTemplateParams{ID: s.templateID, UserId: ownerSub})
	if err != nil {
		t.Fatalf("GetTemplate reload: %v", err)
	}
	if n, err := q.UpdateTemplateIfMatch(ctx, UpdateTemplateIfMatchParams{
		Title: "T3", Subtitle: "s", Description: "d",
		ID: s.templateID, UserId: ownerSub, BaseUpdatedAt: curTmpl.UpdatedAt,
	}); err != nil || n != 1 {
		t.Fatalf("UpdateTemplateIfMatch: n=%d err=%v", n, err)
	}
}

func TestPlanAccessQueries(t *testing.T) {
	db := openTestDB(t)
	q := New(db)
	ctx := context.Background()
	s := seedPlanner(t, q)

	rows, err := q.ListPlanAccessByPlanSync(ctx, ListPlanAccessByPlanSyncParams{
		PlanID: s.planID, UserID: ownerSub, UpdatedSince: "", CursorTs: "", LimitCount: 10,
	})
	if err != nil || len(rows) != 1 || rows[0].User != guestSub {
		t.Fatalf("ListPlanAccessByPlanSync owner: %+v err=%v", rows, err)
	}
	guestRows, err := q.ListPlanAccessByPlanSync(ctx, ListPlanAccessByPlanSyncParams{
		PlanID: s.planID, UserID: guestSub, UpdatedSince: "", CursorTs: "", LimitCount: 10,
	})
	if err != nil || len(guestRows) != 1 {
		t.Fatalf("ListPlanAccessByPlanSync guest: %d err=%v", len(guestRows), err)
	}
	if _, err := q.ListPlanAccessByPlanSync(ctx, ListPlanAccessByPlanSyncParams{
		PlanID: s.planID, UserID: outsiderSub, UpdatedSince: "", CursorTs: "", LimitCount: 10,
	}); err != nil {
		t.Fatalf("ListPlanAccessByPlanSync outsider: %v", err)
	}

	if err := q.RemoveAccess(ctx, guestSub); err != nil {
		t.Fatalf("RemoveAccess: %v", err)
	}
	after, err := q.ListPlanAccess(ctx, ListPlanAccessParams{ID: s.planID, User: ownerSub})
	if err != nil || len(after) != 0 {
		t.Fatalf("ListPlanAccess after RemoveAccess: %d err=%v", len(after), err)
	}
}

func TestIdempotencyQueries(t *testing.T) {
	db := openTestDB(t)
	q := New(db)
	ctx := context.Background()

	if _, err := q.GetIdempotencyRecord(ctx, GetIdempotencyRecordParams{
		UserID: "u1", KeyHash: "abc",
	}); !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("GetIdempotencyRecord missing: %v", err)
	}
	if err := q.InsertIdempotencyResponse(ctx, InsertIdempotencyResponseParams{
		UserID: "u1", KeyHash: "abc", RequestHash: "req1",
		ResponseBody: `{"ok":true}`, StatusCode: 200,
	}); err != nil {
		t.Fatalf("InsertIdempotencyResponse: %v", err)
	}
	row, err := q.GetIdempotencyRecord(ctx, GetIdempotencyRecordParams{UserID: "u1", KeyHash: "abc"})
	if err != nil || row.RequestHash != "req1" || row.StatusCode != 200 {
		t.Fatalf("GetIdempotencyRecord: %+v err=%v", row, err)
	}
}

func TestSyncQueriesAcceptNilOptionalFilters(t *testing.T) {
	db := openTestDB(t)
	q := New(db)
	ctx := context.Background()
	s := seedPlanner(t, q)

	plans, err := q.ListPlansSync(ctx, ListPlansSyncParams{UserID: ownerSub, LimitCount: 10})
	if err != nil || len(plans) != 1 {
		t.Fatalf("ListPlansSync nil filters: %d err=%v", len(plans), err)
	}
	tasks, err := q.ListTasksByPlanSync(ctx, ListTasksByPlanSyncParams{PlanID: s.planID, UserID: ownerSub, LimitCount: 10})
	if err != nil || len(tasks) != 1 {
		t.Fatalf("ListTasksByPlanSync nil filters: %d err=%v", len(tasks), err)
	}
}

func TestSoftDeletesHideRows(t *testing.T) {
	db := openTestDB(t)
	q := New(db)
	ctx := context.Background()
	s := seedPlanner(t, q)

	if err := q.DeleteTask(ctx, DeleteTaskParams{ID: s.taskID, UserId: ownerSub}); err != nil {
		t.Fatalf("DeleteTask: %v", err)
	}
	if _, err := q.GetTask(ctx, GetTaskParams{ID: s.taskID, UserId: ownerSub}); !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("GetTask after delete: %v", err)
	}

	if err := q.DeleteResource(ctx, DeleteResourceParams{ID: s.resourceID, UserId: ownerSub}); err != nil {
		t.Fatalf("DeleteResource: %v", err)
	}
	if _, err := q.GetResource(ctx, GetResourceParams{ID: s.resourceID, UserId: ownerSub}); !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("GetResource after delete: %v", err)
	}

	if err := q.DeleteTemplate(ctx, DeleteTemplateParams{ID: s.templateID, UserId: ownerSub}); err != nil {
		t.Fatalf("DeleteTemplate: %v", err)
	}
	if _, err := q.GetTemplate(ctx, GetTemplateParams{ID: s.templateID, UserId: ownerSub}); !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("GetTemplate after delete: %v", err)
	}

	if err := q.DeletePlan(ctx, DeletePlanParams{ID: s.planID, UserId: ownerSub}); err != nil {
		t.Fatalf("DeletePlan: %v", err)
	}
	if _, err := q.GetPlan(ctx, GetPlanParams{ID: s.planID, UserId: ownerSub}); !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("GetPlan after delete: %v", err)
	}
}
