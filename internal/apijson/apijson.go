package apijson

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Error writes a consistent JSON error envelope for /api/v1.
func Error(c echo.Context, status int, code, message string) error {
	return c.JSON(status, map[string]any{
		"error": map[string]string{
			"code":    code,
			"message": message,
		},
	})
}

// Conflict writes 409 with the server's current entity under error.entity for sync conflicts.
func Conflict(c echo.Context, entity any) error {
	return c.JSON(http.StatusConflict, map[string]any{
		"error": map[string]any{
			"code":    "CONFLICT",
			"message": "The resource changed since your last read; merge server state and retry with a current updated_at (If-Match or base_updated_at).",
			"entity":  entity,
		},
	})
}
