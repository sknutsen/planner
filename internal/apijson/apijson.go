package apijson

import (
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
)

// IncludeDetail controls whether SERVER_ERROR responses include an error.detail field.
// Enabled when API_DEBUG is 1, true, or yes (useful when testing endpoints locally).
var IncludeDetail = apiDebugEnabled()

func apiDebugEnabled() bool {
	v := strings.ToLower(strings.TrimSpace(os.Getenv("API_DEBUG")))
	return v == "1" || v == "true" || v == "yes"
}

// Error writes a consistent JSON error envelope for /api/v1.
func Error(c echo.Context, status int, code, message string) error {
	if status >= http.StatusInternalServerError {
		logRequestError(c, code, message, nil)
	}
	return c.JSON(status, map[string]any{
		"error": map[string]string{
			"code":    code,
			"message": message,
		},
	})
}

// ServerError logs the underlying error and returns a SERVER_ERROR JSON envelope.
// When API_DEBUG is set, the response includes error.detail with err.Error().
func ServerError(c echo.Context, message string, err error) error {
	logRequestError(c, "SERVER_ERROR", message, err)
	errBody := map[string]string{
		"code":    "SERVER_ERROR",
		"message": message,
	}
	if IncludeDetail && err != nil {
		errBody["detail"] = err.Error()
	}
	return c.JSON(http.StatusInternalServerError, map[string]any{"error": errBody})
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

func logRequestError(c echo.Context, code, message string, err error) {
	attrs := []any{
		"method", c.Request().Method,
		"path", c.Request().URL.Path,
		"code", code,
		"message", message,
	}
	if rid := c.Response().Header().Get(echo.HeaderXRequestID); rid != "" {
		attrs = append(attrs, "request_id", rid)
	}
	if uid, ok := c.Get("api_user_sub").(string); ok && uid != "" {
		attrs = append(attrs, "user", uid)
	}
	if err != nil {
		attrs = append(attrs, "err", err)
	}
	slog.Error("api error", attrs...)
}
