package view

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/sknutsen/planner/models"
)

func TestTaskFromTemplate_RendersTitleAndSubtitle(t *testing.T) {
	var buf bytes.Buffer
	err := TaskFromTemplate(models.ClientState{SelectedPlanId: 3}, models.Task{Date: "2026-08-26"}, []models.Template{
		{Id: 11, Title: "Standup notes", Subtitle: "Daily 15-minute team sync"},
		{Id: 12, Title: "This is a very long template title that should wrap onto multiple lines in the selector", Subtitle: "An equally long subtitle that also needs room to wrap under the title"},
		{Id: 13, Title: "Title only template", Subtitle: ""},
		{Id: 14, Title: "", Subtitle: "Subtitle only template"},
		{Id: 15, Title: "   ", Subtitle: ""},
	}).Render(context.Background(), &buf)
	if err != nil {
		t.Fatal(err)
	}
	html := buf.String()

	mustContain := []string{
		`name="template"`,
		`id="template-selector"`,
		`Select a template`,
		"Standup notes",
		"Daily 15-minute team sync",
		"This is a very long template title that should wrap onto multiple lines in the selector",
		"An equally long subtitle that also needs room to wrap under the title",
		"Title only template",
		"Subtitle only template",
		"Untitled template",
		`data-template-id="11"`,
		`data-template-id="15"`,
		"template-selector__title",
		"template-selector__subtitle",
		"template-selector__untitled",
	}
	for _, s := range mustContain {
		if !strings.Contains(html, s) {
			t.Errorf("expected %q in rendered HTML", s)
		}
	}

	titleOnly := optionHTML(t, html, "13")
	if !strings.Contains(titleOnly, "Title only template") {
		t.Errorf("title-only option missing title: %s", titleOnly)
	}
	if strings.Contains(titleOnly, "template-selector__subtitle") || strings.Contains(titleOnly, "Untitled template") {
		t.Errorf("title-only option should omit subtitle and untitled fallback: %s", titleOnly)
	}

	subtitleOnly := optionHTML(t, html, "14")
	if !strings.Contains(subtitleOnly, "Subtitle only template") {
		t.Errorf("subtitle-only option missing subtitle: %s", subtitleOnly)
	}
	if strings.Contains(subtitleOnly, "template-selector__title") || strings.Contains(subtitleOnly, "Untitled template") {
		t.Errorf("subtitle-only option should omit title and untitled fallback: %s", subtitleOnly)
	}

	untitled := optionHTML(t, html, "15")
	if !strings.Contains(untitled, "Untitled template") {
		t.Errorf("blank option missing untitled fallback: %s", untitled)
	}
	if strings.Contains(untitled, "template-selector__title") || strings.Contains(untitled, "template-selector__subtitle") {
		t.Errorf("blank option should not render empty title/subtitle spans: %s", untitled)
	}
}

func TestTaskFromTemplate_EmptyList(t *testing.T) {
	var buf bytes.Buffer
	err := TaskFromTemplate(models.ClientState{}, models.Task{Date: "2026-08-26"}, nil).Render(context.Background(), &buf)
	if err != nil {
		t.Fatal(err)
	}
	html := buf.String()
	if !strings.Contains(html, "No templates") {
		t.Fatalf("expected empty-state copy, got %s", html)
	}
	if strings.Contains(html, `role="option"`) {
		t.Fatal("empty list should not render options")
	}
}

func optionHTML(t *testing.T, html, id string) string {
	t.Helper()
	marker := `data-template-id="` + id + `"`
	start := strings.Index(html, marker)
	if start < 0 {
		t.Fatalf("option %s not found", id)
	}
	rest := html[start:]
	next := strings.Index(rest[len(marker):], `data-template-id="`)
	if next >= 0 {
		return rest[:len(marker)+next]
	}
	return rest
}
