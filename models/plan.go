package models

type Plan struct {
	Id   int
	Name string
}

type UpdatePlanRequest struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
