package auth

import (
	log "github.com/akakou/ra-webs/service"
	"github.com/labstack/echo/v4"
)

func Authenticate(l *log.Log, c echo.Context) error {
	authorization := c.Request().Header["Authorization"][0]
	token := authorization[len("Bearer "):]

	if token != l.Token {
		return echo.ErrUnauthorized
	}

	return nil
}
