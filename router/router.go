package router

import (
	"log/slog"
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
	setupHTTPErrorHandler(e)

	e.Use(middleware.RequestID())
	e.Use(middlewareRecover())
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogMethod:   true,
		LogLatency:  true,
		LogError:    true,
		LogRemoteIP: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			attrs := []any{
				"method", v.Method,
				"path", v.URI,
				"status", v.Status,
				"latency", v.Latency.String(),
				"remote_ip", v.RemoteIP,
			}
			if rid := c.Response().Header().Get(echo.HeaderXRequestID); rid != "" {
				attrs = append(attrs, "request_id", rid)
			}
			if v.Error != nil {
				attrs = append(attrs, "err", v.Error)
			}

			switch {
			case v.Status >= http.StatusInternalServerError:
				slog.Error("request", attrs...)
			case v.Status >= http.StatusBadRequest:
				slog.Warn("request", attrs...)
			default:
				slog.Info("request", attrs...)
			}
			return nil
		},
	}))

	handler.RegisterSwagger(e)

	e.Static(routes.Assets, "assets")

	api := e.Group("/api/v1")
	api.Use(mw.BearerJWT(&h.Authenticator, h.AuthConfig.APIAudience))
	api.Use(mw.IdempotencyBuffer(h.DB))
	h.RegisterAPIV1(api)

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

	e.GET(routes.Resources, h.Resources, mw.IsAuthenticated)
	e.GET(routes.ResourcesPlan, h.Resources, mw.IsAuthenticated)
	e.GET(routes.ResourcesPlanResources, h.ListAllResources, mw.IsAuthenticated)

	e.GET(routes.Resource, h.Resource, mw.IsAuthenticated)
	e.GET(routes.ResourceCreate, h.CreateResource, mw.IsAuthenticated)
	e.GET(routes.ResourceDelete, h.DeleteResource, mw.IsAuthenticated)
	e.GET(routes.ResourceEdit, h.EditResource, mw.IsAuthenticated)
	e.POST(routes.ResourceUpdate, h.UpdateResource, mw.IsAuthenticated)

	e.GET(routes.Task, h.Task, mw.IsAuthenticated)
	e.GET(routes.TaskComplete, h.ToggleIsCompleteTask, mw.IsAuthenticated)
	e.POST(routes.TaskCopy, h.CopyTask, mw.IsAuthenticated)
	e.GET(routes.TaskCreate, h.CreateTask, mw.IsAuthenticated)
	e.GET(routes.TaskCreateTemplate, h.CreateTaskFromTemplate, mw.IsAuthenticated)
	e.GET(routes.TaskDelete, h.DeleteTask, mw.IsAuthenticated)
	e.GET(routes.TaskEdit, h.EditTask, mw.IsAuthenticated)
	e.POST(routes.TaskUpdate, h.UpdateTask, mw.IsAuthenticated)

	e.GET(routes.Templates, h.Templates, mw.IsAuthenticated)
	e.GET(routes.TemplatesPlan, h.Templates, mw.IsAuthenticated)
	e.GET(routes.TemplatesPlanTemplates, h.ListAllTemplates, mw.IsAuthenticated)

	e.GET(routes.Template, h.Template, mw.IsAuthenticated)
	e.GET(routes.TemplateCreate, h.CreateTemplate, mw.IsAuthenticated)
	e.GET(routes.TemplateDelete, h.DeleteTemplate, mw.IsAuthenticated)
	e.GET(routes.TemplateEdit, h.EditTemplate, mw.IsAuthenticated)
	e.POST(routes.TemplateUpdate, h.UpdateTemplate, mw.IsAuthenticated)
	e.POST(routes.TemplateFromTask, h.TemplateFromTask, mw.IsAuthenticated)

	e.GET(routes.User, h.User, mw.IsAuthenticated)

	e.GET(routes.Week, h.Week, mw.IsAuthenticated)
	e.GET(routes.WeekPlan, h.Week, mw.IsAuthenticated)
	e.GET(routes.WeekPlanWeek, h.Week, mw.IsAuthenticated)

	e.GET("/favicon.ico", func(c echo.Context) error {
		return c.NoContent(http.StatusNoContent)
	})
}
