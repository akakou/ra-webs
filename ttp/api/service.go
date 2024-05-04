package api

import (
	"net/http"

	goutils "github.com/akakou/go-utils"
	ttpcore "github.com/akakou/ra_webs/ttp/core"
	"github.com/labstack/echo/v4"
)

var PostServiceByAdmin = goutils.EchoRoute[ttpcore.TTP]{
	Method: goutils.POST,
	Path:   API_ROOT + "/service",
	F: func(ttp *ttpcore.TTP) goutils.EchoRouteFunc {
		return func(c echo.Context) error {
			err := authenticateAdmin(ttp, c)

			if err != nil {
				return c.String(http.StatusUnauthorized, "token is invalid")
			}

			token, err := goutils.RandomHex(ttpcore.RANDOM_SIZE)

			if err != nil {
				return err
			}

			service, err := ttp.DB.Client.Service.
				Create().
				SetName("").
				SetToken(token).
				SetIsActive(true).
				Save(*ttp.DB.Ctx)

			if err != nil {
				return err
			}

			return c.String(http.StatusOK, service.Token)
		}
	},
}
