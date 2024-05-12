package handler

import (
	"log"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
)

func (handler *Handler) Logout(c echo.Context) error {
	logoutUrl, err := url.Parse(handler.AuthConfig.Domain + "/v2/logout")
	if err != nil {
		return err
	}

	log.Println(logoutUrl)

	scheme := "http"
	if c.Request().TLS != nil {
		scheme = "https"
	}

	returnTo, err := url.Parse(scheme + "://" + string(c.Request().Host))
	if err != nil {
		return err
	}

	log.Println(returnTo)

	parameters := url.Values{}
	parameters.Add("returnTo", returnTo.String())
	parameters.Add("client_id", handler.AuthConfig.ClientId)
	logoutUrl.RawQuery = parameters.Encode()

	return c.Redirect(http.StatusTemporaryRedirect, logoutUrl.String())
}
