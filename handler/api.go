package handler

import (
	"errors"

	"github.com/labstack/echo/v4"
	mw "github.com/sknutsen/planner/middleware"
)

func apiUserID(c echo.Context) (string, error) {
	v := c.Get(mw.CtxAPIUserID)
	s, ok := v.(string)
	if !ok || s == "" {
		return "", errors.New("api user missing")
	}
	return s, nil
}
