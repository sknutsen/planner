package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func (handler *Handler) Callback(c echo.Context) error {
	sess, _ := session.Get("session", c)

	if c.QueryParam("state") != sess.Values["state"] {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Bad request. session state: %s\n", sess.Values["state"]))
	}

	code := c.QueryParam("code")

	token, err := handler.Authenticator.Exchange(c.Request().Context(), code)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to exchange an authorization code for a token.")
	}

	idToken, err := handler.Authenticator.VerifyIDToken(c.Request().Context(), token)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to verify ID Token.")
	}

	println("getting profile")
	var profile map[string]interface{}
	if err := idToken.Claims(&profile); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	println("setting token")
	sess.Values["access_token"] = token.AccessToken
	sess.Values["profile"] = profile
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.Redirect(http.StatusTemporaryRedirect, "/")
}
