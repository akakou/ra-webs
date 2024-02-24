package ttp

import (
	"net/http"

	goutils "github.com/akakou/go-utils"
	"github.com/labstack/echo/v4"
)

var postServiceByAdmin = echoRoute{
	method: POST,
	path:   "/service",
	f: func(auditor *Auditor) echoRouteFunc {
		return func(c echo.Context) error {
			err := authenticateAdmin(auditor, c)

			if err != nil {
				return c.String(http.StatusUnauthorized, "token is invalid")
			}

			token, err := goutils.RandomHex(32)

			if err != nil {
				c.Error(err)
				return err
			}

			service, err := auditor.db.Client.Service.Create().SetName("").SetToken(token).SetHasActivated(true).Save(*auditor.db.Ctx)

			if err != nil {
				c.Error(err)
				return err
			}

			return c.String(http.StatusOK, service.Token)
		}
	},
}
