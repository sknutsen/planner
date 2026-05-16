package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sknutsen/planner/docs"
)

// RegisterSwagger serves the OpenAPI document and Swagger UI at /swagger (no authentication).
func RegisterSwagger(e *echo.Echo) {
	e.GET("/swagger/openapi.yaml", func(c echo.Context) error {
		return c.Blob(http.StatusOK, "application/yaml", docs.OpenAPIYAML)
	})
	e.GET("/swagger", func(c echo.Context) error {
		return c.Blob(http.StatusOK, "text/html; charset=utf-8", docs.SwaggerUIPage)
	})
}
