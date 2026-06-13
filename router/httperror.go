package router

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/sknutsen/planner/internal/apijson"
)

func isAPIRequest(c echo.Context) bool {
	return strings.HasPrefix(c.Request().URL.Path, "/api/v1")
}

// setupHTTPErrorHandler logs unhandled errors and returns consistent JSON for /api/v1.
func setupHTTPErrorHandler(e *echo.Echo) {
	defaultHandler := e.HTTPErrorHandler
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		if isAPIRequest(c) {
			var he *echo.HTTPError
			if errors.As(err, &he) {
				code := httpStatusToAPICode(he.Code)
				msg := httpErrorMessage(he)
				if he.Code >= http.StatusInternalServerError {
					apijson.ServerError(c, msg, err)
					return
				}
				_ = apijson.Error(c, he.Code, code, msg)
				return
			}
			apijson.ServerError(c, "An unexpected error occurred.", err)
			return
		}

		defaultHandler(err, c)
	}
}

func setupRecover(e *echo.Echo) {
	e.Use(middlewareRecover())
}

func middlewareRecover() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			defer func() {
				if r := recover(); r != nil {
					err = recoverToError(r)
					attrs := []any{
						"method", c.Request().Method,
						"path", c.Request().URL.Path,
						"panic", r,
					}
					if rid := c.Response().Header().Get(echo.HeaderXRequestID); rid != "" {
						attrs = append(attrs, "request_id", rid)
					}
					slog.Error("panic recovered", attrs...)

					if isAPIRequest(c) && !c.Response().Committed {
						_ = apijson.ServerError(c, "An unexpected error occurred.", err)
						err = nil
					}
				}
			}()
			return next(c)
		}
	}
}

func recoverToError(r any) error {
	switch v := r.(type) {
	case error:
		return v
	case string:
		return errors.New(v)
	default:
		return errors.New("panic")
	}
}

func httpStatusToAPICode(status int) string {
	switch status {
	case http.StatusBadRequest:
		return "BAD_REQUEST"
	case http.StatusUnauthorized:
		return "UNAUTHORIZED"
	case http.StatusNotFound:
		return "NOT_FOUND"
	case http.StatusConflict:
		return "CONFLICT"
	case http.StatusServiceUnavailable:
		return "MISCONFIGURED"
	default:
		if status >= http.StatusInternalServerError {
			return "SERVER_ERROR"
		}
		return "BAD_REQUEST"
	}
}

func httpErrorMessage(he *echo.HTTPError) string {
	switch msg := he.Message.(type) {
	case string:
		if msg != "" {
			return msg
		}
	case error:
		if msg != nil {
			return msg.Error()
		}
	}
	return http.StatusText(he.Code)
}
