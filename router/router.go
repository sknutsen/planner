package router

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sknutsen/planner/handler"
	mw "github.com/sknutsen/planner/middleware"
	"github.com/sknutsen/planner/routes"
)

func Setup(e *echo.Echo, h *handler.Handler) {
	// e.Use(middleware.Logger())
	// e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
	// 	Format: "time=${time_rfc3339} method=${method}, uri=${uri}, status=${status}\n",
	// }))
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus: true,
		LogURI:    true,
		LogMethod: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			fmt.Printf("%v %v %v status: %v error: %v\n", v.StartTime.Format(time.RFC822), v.Method, v.URI, v.Status, v.Error)
			return nil
		},
	}))

	e.Static(routes.Assets, "assets")

	e.Use(session.Middleware(sessions.NewCookieStore([]byte(h.AuthConfig.ClientSecret))))

	e.GET(routes.Index, h.Index, mw.IsAuthenticated)

	e.GET(routes.Callback, h.Callback)

	e.GET(routes.ComponentsModal, h.Modal, mw.IsAuthenticated)

	e.GET(routes.Day, h.Day, mw.IsAuthenticated)
	e.GET(routes.Daytasks, h.DayTasks, mw.IsAuthenticated)

	e.GET(routes.History, h.History, mw.IsAuthenticated)
	e.GET(routes.HistoryPlan, h.History, mw.IsAuthenticated)
	e.GET(routes.HistoryPlanTasks, h.ListAllTasks, mw.IsAuthenticated)

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

	e.GET(routes.User, h.User, mw.IsAuthenticated)

	e.GET(routes.Week, h.Week, mw.IsAuthenticated)
	e.GET(routes.WeekPlan, h.Week, mw.IsAuthenticated)
	e.GET(routes.WeekPlanWeek, h.Week, mw.IsAuthenticated)

	e.GET("/favicon.ico", func(c echo.Context) error {
		return c.NoContent(http.StatusNoContent)
	})
}
