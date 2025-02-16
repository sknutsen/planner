package models

import "github.com/sknutsen/planner/database"

type Template struct {
	Id          int
	Title       string
	Subtitle    string
	Description string
}

type UpdateTemplateRequest struct {
	Id          string `json:"id"`
	PlanId      string `json:"plan_id"`
	Title       string `json:"title"`
	Subtitle    string `json:"subtitle"`
	Description string `json:"description"`
}

type PlanTemplatesResponse struct {
	Templates []Template
}

func TemplatesFromDBModels(m []database.Template) []Template {
	templates := []Template{}

	for _, v := range m {
		templates = append(templates, TemplateFromDBModel(v))
	}

	return templates
}

func TemplateFromDBModel(m database.Template) Template {
	return Template{
		Id:          int(m.ID),
		Title:       m.Title,
		Subtitle:    m.Subtitle.(string),
		Description: m.Description.(string),
	}
}
