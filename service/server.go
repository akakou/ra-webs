package service

import (
	"net/http"

	"github.com/akakou/ra_webs/core"
	"github.com/labstack/echo/v4"
)

func (service *Service) ServDomainAuth(token, domain string, e *echo.Echo) *echo.Echo {
	e.GET(core.DOMAIN_AUTH_PATH, func(c echo.Context) error {
		return c.String(http.StatusOK, token)
	})

	go e.Start(core.ServicePort)
	return e
}
