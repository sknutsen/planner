package router

import (
	"github.com/labstack/echo/v4"
	"github.com/sknutsen/planner/handler"
	"github.com/sknutsen/planner/routes"
)

func Setup(e *echo.Echo, h *handler.Handler) {
	e.Static(routes.Assets, "assets")

	e.GET(routes.Index, h.Index)
	e.GET(routes.IndexWeek, h.Index)

	e.GET(routes.ComponentsModal, h.Modal)

	e.GET(routes.Day, h.Day)
	e.GET(routes.Daytasks, h.DayTasks)

	e.GET(routes.Task, h.Task)
	e.GET(routes.TaskCreate, h.CreateTask)
	e.GET(routes.TaskDelete, h.DeleteTask)
	e.GET(routes.TaskEdit, h.EditTask)
	e.POST(routes.TaskUpdate, h.UpdateTask)
}
