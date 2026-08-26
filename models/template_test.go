package models

import (
	"testing"

	"github.com/sknutsen/planner/database"
)

func TestTemplateFromDBModel(t *testing.T) {
	m := database.Template{
		ID:          7,
		PlanID:      2,
		Title:       "Standup",
		Subtitle:    "Daily sync",
		Description: "Notes",
	}
	got := TemplateFromDBModel(m)
	if got.Id != 7 || got.Title != "Standup" || got.Subtitle != "Daily sync" || got.Description != "Notes" {
		t.Fatalf("got %+v", got)
	}
}

func TestTemplateFromDBModel_NilOptionalStrings(t *testing.T) {
	m := database.Template{
		ID:          1,
		PlanID:      1,
		Title:       "",
		Subtitle:    nil,
		Description: nil,
	}
	got := TemplateFromDBModel(m)
	if got.Title != "" || got.Subtitle != "" || got.Description != "" {
		t.Fatalf("want empty strings, got %+v", got)
	}
}

func TestTemplatesFromDBModels(t *testing.T) {
	out := TemplatesFromDBModels([]database.Template{
		{ID: 1, Title: "a", Subtitle: nil, Description: nil},
		{ID: 2, Title: "b", Subtitle: "sub", Description: "desc"},
	})
	if len(out) != 2 || out[0].Id != 1 || out[0].Subtitle != "" || out[1].Subtitle != "sub" {
		t.Fatalf("got %+v", out)
	}
}
