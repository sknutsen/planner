package models

import (
	"testing"

	"github.com/sknutsen/planner/database"
)

func TestTaskFromDBModel(t *testing.T) {
	m := database.Task{
		ID:          10,
		PlanID:      2,
		Date:        "2025-03-01T00:00:00Z",
		Title:       "Title",
		Subtitle:    "Sub",
		Description: "Desc",
		IsComplete:  1,
	}
	got := TaskFromDBModel(m)
	if got.Id != 10 || got.Date != "2025-03-01" || got.Title != "Title" {
		t.Fatalf("basic fields: %+v", got)
	}
	if got.Subtitle != "Sub" || got.Description != "Desc" {
		t.Fatalf("text fields: %+v", got)
	}
	if !got.IsComplete {
		t.Fatal("expected IsComplete true")
	}
}

func TestTaskFromDBModel_NotComplete(t *testing.T) {
	m := database.Task{
		ID:          1,
		PlanID:      1,
		Date:        "2025-01-01",
		Title:       "T",
		Subtitle:    "",
		Description: "",
		IsComplete:  0,
	}
	got := TaskFromDBModel(m)
	if got.IsComplete {
		t.Fatal("expected IsComplete false")
	}
}

func TestTaskFromDBModel_NilOptionalStrings(t *testing.T) {
	m := database.Task{
		ID:          1,
		PlanID:      1,
		Date:        "2025-01-01",
		Title:       "T",
		Subtitle:    nil,
		Description: nil,
		IsComplete:  0,
	}
	got := TaskFromDBModel(m)
	if got.Subtitle != "" || got.Description != "" {
		t.Fatalf("want empty strings, got %+v", got)
	}
}

func TestTasksFromDBModels(t *testing.T) {
	out := TasksFromDBModels([]database.Task{
		{ID: 1, PlanID: 1, Date: "2025-01-01", Title: "a", Subtitle: "", Description: "", IsComplete: 0},
		{ID: 2, PlanID: 1, Date: "2025-01-02", Title: "b", Subtitle: "", Description: "", IsComplete: 0},
	})
	if len(out) != 2 || out[0].Id != 1 || out[1].Title != "b" {
		t.Fatalf("got %+v", out)
	}
}
