package models

import "github.com/sknutsen/planner/database"

type Task struct {
	Id          int
	Date        string
	Title       string
	Subtitle    string
	Description string
	IsComplete  bool
}

type UpdateTaskRequest struct {
	Id          string `json:"id"`
	PlanId      string `json:"plan_id"`
	Date        string `json:"date"`
	Title       string `json:"title"`
	Subtitle    string `json:"subtitle"`
	Description string `json:"description"`
}

func TasksFromDBModels(m []database.Task) []Task {
	tasks := []Task{}

	for _, v := range m {
		tasks = append(tasks, TaskFromDBModel(v))
	}

	return tasks
}

func TaskFromDBModel(m database.Task) Task {
	return Task{
		Id:          int(m.ID),
		Date:        m.Date,
		Title:       m.Title,
		Subtitle:    m.Subtitle.(string),
		Description: m.Description.(string),
		IsComplete:  m.IsComplete != 0,
	}
}
