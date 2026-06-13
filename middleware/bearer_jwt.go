package middleware

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/labstack/echo/v4"
	"github.com/sknutsen/planner/auth"
	"github.com/sknutsen/planner/internal/apijson"
)

// CtxAPIUserID is the Echo context key for the Auth0 subject (same as session UserProfile.UserId).
const CtxAPIUserID = "api_user_sub"

// BearerJWT validates Authorization: Bearer access tokens against the Auth0 OIDC provider.
// audience must be the Auth0 API identifier (AUTH0_API_AUDIENCE); this is the access token aud, not the SPA client_id.
func BearerJWT(a *auth.Authenticator, audience string) echo.MiddlewareFunc {
	if audience == "" {
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				return apijson.Error(c, http.StatusServiceUnavailable, "MISCONFIGURED", "AUTH0_API_AUDIENCE is not set")
		}
		}
	}
	verifier := a.Verifier(&oidc.Config{ClientID: audience})
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			raw := c.Request().Header.Get("Authorization")
			if raw == "" || !strings.HasPrefix(strings.ToLower(raw), "bearer ") {
				return apijson.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Missing bearer token")
			}
			token := strings.TrimSpace(raw[len("Bearer "):])
			idt, err := verifier.Verify(c.Request().Context(), token)
			if err != nil {
				slog.Debug("jwt verification failed", "path", c.Request().URL.Path, "err", err)
				return apijson.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid token")
			}
			var claims struct {
				Sub string `json:"sub"`
			}
			if err := idt.Claims(&claims); err != nil || claims.Sub == "" {
				slog.Debug("jwt claims invalid", "path", c.Request().URL.Path, "err", err)
				return apijson.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid claims")
			}
			c.Set(CtxAPIUserID, claims.Sub)
			return next(c)
		}
	}
}
