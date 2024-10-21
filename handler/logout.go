package handler

import (
	"net/http"
	"net/url"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func (handler *Handler) Logout(c echo.Context) error {
	logoutUrl, err := url.Parse(handler.AuthConfig.Domain + "/v2/logout")
	if err != nil {
		return err
	}

	scheme := "http"
	if c.Request().TLS != nil {
		scheme = "https"
	}

	returnTo, err := url.Parse(scheme + "://" + string(c.Request().Host))
	if err != nil {
		return err
	}

	parameters := url.Values{}
	parameters.Add("returnTo", returnTo.String())
	parameters.Add("client_id", handler.AuthConfig.ClientId)
	logoutUrl.RawQuery = parameters.Encode()

	sess, _ := session.Get("session", c)

	println("setting token")
	sess.Values["access_token"] = nil
	sess.Values["profile"] = nil
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.Redirect(http.StatusTemporaryRedirect, logoutUrl.String())
}
