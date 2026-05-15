package routes

import "fmt"

func HistoryPlanTasksPath(planID int) string {
	return fmt.Sprintf("/history/%d/tasks", planID)
}

func PlanTemplatesListPath(planID int) string {
	return fmt.Sprintf("/templates/%d/templates", planID)
}

func PlanResourcesListPath(planID int) string {
	return fmt.Sprintf("/resources/%d/resources", planID)
}

func PlanEditPath(planID int) string {
	return fmt.Sprintf("/plan/%d/edit", planID)
}

func TemplateCreatePath(planID int) string {
	return fmt.Sprintf("/templates/%d/create", planID)
}

func ResourceCreatePath(planID int) string {
	return fmt.Sprintf("/resources/%d/create", planID)
}

func TemplatePath(id int) string {
	return fmt.Sprintf("/template/%d", id)
}

func TemplateDeletePath(id int) string {
	return fmt.Sprintf("/template/%d/delete", id)
}

func ResourcePath(id int) string {
	return fmt.Sprintf("/resource/%d", id)
}

func ResourceDeletePath(id int) string {
	return fmt.Sprintf("/resource/%d/delete", id)
}
