package handler

import (
	"errors"
	"fmt"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/sknutsen/planner/models"
)

var errNoSessionProfile = errors.New("session profile missing or invalid")

// userProfileFromContext reads the signed session and maps the OIDC profile into
// models.UserProfile. It never type-asserts blindly on the profile map.
func userProfileFromContext(c echo.Context) (models.UserProfile, error) {
	sess, err := session.Get("session", c)
	if err != nil {
		return models.UserProfile{}, fmt.Errorf("session: %w", err)
	}
	raw, ok := sess.Values["profile"]
	if !ok || raw == nil {
		return models.UserProfile{}, errNoSessionProfile
	}
	m, ok := raw.(map[string]interface{})
	if !ok {
		return models.UserProfile{}, errNoSessionProfile
	}
	return models.GetUserProfile(m), nil
}
