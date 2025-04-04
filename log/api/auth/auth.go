package auth

import (
	"github.com/akakou/ra-webs/log/core"
	"github.com/labstack/echo/v4"
)

func Authenticate(log *core.Log, c echo.Context) error {
	authorization := c.Request().Header["Authorization"][0]
	token := authorization[len("Bearer "):]

	if token != log.Token {
		return echo.ErrUnauthorized
	}

	return nil
}
