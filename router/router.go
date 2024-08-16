package router

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sknutsen/planner/handler"
	mw "github.com/sknutsen/planner/middleware"
	"github.com/sknutsen/planner/routes"
)

func Setup(e *echo.Echo, h *handler.Handler) {
	e.Use(middleware.Logger())

	e.Static(routes.Assets, "assets")

	e.Use(session.Middleware(sessions.NewCookieStore([]byte(h.AuthConfig.ClientSecret))))

	e.GET(routes.Index, h.Index, mw.IsAuthenticated)
	e.GET(routes.IndexPlan, h.Index, mw.IsAuthenticated)
	e.GET(routes.IndexPlanWeek, h.Index, mw.IsAuthenticated)

	e.GET(routes.Callback, h.Callback)

	e.GET(routes.ComponentsModal, h.Modal, mw.IsAuthenticated)

	e.GET(routes.Day, h.Day, mw.IsAuthenticated)
	e.GET(routes.Daytasks, h.DayTasks, mw.IsAuthenticated)

	e.GET(routes.Login, h.Login)
	e.GET(routes.Logout, h.Logout)

	e.GET(routes.Plan, h.Plan, mw.IsAuthenticated)
	e.GET(routes.PlanCreate, h.CreatePlan, mw.IsAuthenticated)
	e.GET(routes.PlanDelete, h.DeletePlan, mw.IsAuthenticated)
	e.GET(routes.PlanEdit, h.EditPlan, mw.IsAuthenticated)
	e.POST(routes.PlanUpdate, h.UpdatePlan, mw.IsAuthenticated)

	e.GET(routes.Task, h.Task, mw.IsAuthenticated)
	e.GET(routes.TaskComplete, h.ToggleIsCompleteTask, mw.IsAuthenticated)
	e.POST(routes.TaskCopy, h.CopyTask, mw.IsAuthenticated)
	e.GET(routes.TaskCreate, h.CreateTask, mw.IsAuthenticated)
	e.GET(routes.TaskDelete, h.DeleteTask, mw.IsAuthenticated)
	e.GET(routes.TaskEdit, h.EditTask, mw.IsAuthenticated)
	e.POST(routes.TaskUpdate, h.UpdateTask, mw.IsAuthenticated)

	e.GET("/favicon.ico", func(c echo.Context) error {
		return c.NoContent(http.StatusNoContent)
	})
}
