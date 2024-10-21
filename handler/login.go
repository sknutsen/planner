package handler

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func (handler *Handler) Login(c echo.Context) error {
	state, err := generateRandomState()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	sess, _ := session.Get("session", c)
	sess.Values["state"] = state
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.Redirect(http.StatusTemporaryRedirect, handler.Authenticator.AuthCodeURL(state))
}

func generateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	state := base64.StdEncoding.EncodeToString(b)

	return state, nil
}
