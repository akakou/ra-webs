package api

import (
	"fmt"

	"github.com/akakou/ra_webs/monitor/serv"
	"github.com/labstack/echo/v4"
)

const (
	ERROR_AUTHENTICATE_SERVICE      = "failed to authenticate service"
	ERROR_SERVICE_NOT_ACTIVE        = "service is not active"
	ERROR_AUTHENTICATE_ADMIN        = "failed to authenticate admin"
	ERROR_ACCESS_DOMAIN_AUTH_TARGET = "failed to access domain auth target"
	ERROR_DOMAIN_AUTH_INVALID       = "domain auth token is invalid"
	ERROR_QUOTE_INVALID1            = "quote is invalid (debug)"
	ERROR_QUOTE_INVALID2            = "quote is invalid (up-to-date)"
	ERROR_QUOTE_INVALID3            = "quote is invalid (unique)"
)

var SCHEME = "https"

func authenticateAdmin(serv *serv.MonitorServer, c echo.Context) error {
	authorization := c.Request().Header["Authorization"][0]
	token := authorization[len("Bearer "):]

	if token != serv.AdminToken {
		return fmt.Errorf(ERROR_AUTHENTICATE_ADMIN)
	}

	return nil
}
