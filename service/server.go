package service

import (
	"net/http"

	goutils "github.com/akakou/go-utils"
	"github.com/akakou/ra_webs/core"
	"github.com/labstack/echo/v4"
)

func (service *Service) ServAuthDomain() {
	e := echo.New()

	e.GET(core.DOMAIN_AUTH_PATH, func(c echo.Context) error {
		nonce, _ := goutils.RandomHex(64)
		h := core.DomainToken(service.Token, nonce)
		return c.String(http.StatusOK, h)
	})

	e.Start(":8081")
}
