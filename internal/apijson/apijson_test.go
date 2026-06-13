package apijson_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/sknutsen/planner/internal/apijson"
)

func TestServerErrorWithoutDetail(t *testing.T) {
	apijson.IncludeDetail = false

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/plans", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := apijson.ServerError(c, "Could not list plans.", errors.New("db timeout"))
	if err != nil {
		t.Fatalf("ServerError returned error: %v", err)
	}
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}

	var body map[string]map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if body["error"]["code"] != "SERVER_ERROR" {
		t.Fatalf("code = %q", body["error"]["code"])
	}
	if body["error"]["message"] != "Could not list plans." {
		t.Fatalf("message = %q", body["error"]["message"])
	}
	if _, ok := body["error"]["detail"]; ok {
		t.Fatal("detail should be omitted when API_DEBUG is off")
	}
}

func TestServerErrorWithDetail(t *testing.T) {
	apijson.IncludeDetail = true
	t.Cleanup(func() { apijson.IncludeDetail = apijsonDebugFromEnv() })

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/plans", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	root := errors.New("db timeout")
	if err := apijson.ServerError(c, "Could not list plans.", root); err != nil {
		t.Fatalf("ServerError returned error: %v", err)
	}

	var body map[string]map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if body["error"]["detail"] != "db timeout" {
		t.Fatalf("detail = %q, want %q", body["error"]["detail"], "db timeout")
	}
}

func apijsonDebugFromEnv() bool {
	v := os.Getenv("API_DEBUG")
	return v == "1" || v == "true" || v == "yes"
}
