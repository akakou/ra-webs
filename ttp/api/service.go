package api

import (
	"net/http"

	goutils "github.com/akakou/go-utils"
	ttpcore "github.com/akakou/ra_webs/ttp/core"
	"github.com/labstack/echo/v4"
)

var postServiceByAdmin = goutils.EchoRoute[ttpcore.TTP]{
	Method: goutils.POST,
	Path:   "/service",
	F: func(ttp *ttpcore.TTP) goutils.EchoRouteFunc {
		return func(c echo.Context) error {
			err := authenticateAdmin(ttp, c)

			if err != nil {
				return c.String(http.StatusUnauthorized, "token is invalid")
			}

			token, err := goutils.RandomHex(32)

			if err != nil {
				c.Error(err)
				return err
			}

			service, err := ttp.DB.Client.Service.Create().SetName("").SetToken(token).SetHasActivated(true).Save(*ttp.DB.Ctx)

			if err != nil {
				c.Error(err)
				return err
			}

			return c.String(http.StatusOK, service.Token)
		}
	},
}
