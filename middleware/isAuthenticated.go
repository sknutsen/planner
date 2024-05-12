package middleware

import (
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/sknutsen/planner/routes"
)

func IsAuthenticated(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get("session", c)
		if err != nil {
			println(err)
			return err
		}

		profile := sess.Values["profile"]

		if profile == nil {
			println("profile is nil")
			return c.Redirect(http.StatusTemporaryRedirect, routes.Login)
		}
		return next(c)
	}
}
