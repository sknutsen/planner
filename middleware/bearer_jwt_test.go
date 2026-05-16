package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/sknutsen/planner/auth"
)

func TestBearerJWT_RequiresAudience(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	next := func(c echo.Context) error {
		return c.NoContent(http.StatusNoContent)
	}
	h := BearerJWT(&auth.Authenticator{}, "")(next)
	if err := h(c); err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503 when audience empty, got %d", rec.Code)
	}
}
