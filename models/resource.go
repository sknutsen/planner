package models

import "github.com/sknutsen/planner/database"

type Resource struct {
	Id      int
	Title   string
	Type    int
	Content string
}

type UpdateResourceRequest struct {
	Id      string `json:"id"`
	PlanId  string `json:"plan_id"`
	Title   string `json:"title"`
	Type    string `json:"type"`
	Content string `json:"content"`
}

type PlanResourcesResponse struct {
	Resources []Resource
}

func ResourcesFromDBModels(m []database.Resource) []Resource {
	tasks := []Resource{}

	for _, v := range m {
		tasks = append(tasks, ResourceFromDBModel(v))
	}

	return tasks
}

func ResourceFromDBModel(m database.Resource) Resource {
	return Resource{
		Id:      int(m.ID),
		Title:   m.Title,
		Type:    int(m.ResourceType),
		Content: m.Content.(string),
	}
}
