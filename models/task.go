package models

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
