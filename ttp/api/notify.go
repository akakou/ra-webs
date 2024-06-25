package api

import (
	"net/http"

	goutils "github.com/akakou/go-utils"
	rawebscore "github.com/akakou/ra_webs/core"
	"github.com/akakou/ra_webs/ttp/core"
	"github.com/akakou/ra_webs/ttp/ent"
	"github.com/akakou/ra_webs/ttp/ent/taserver"
	"github.com/labstack/echo/v4"
)

var PostNotifyApi = goutils.EchoRoute[core.TTP]{
	Method: goutils.POST,
	Path:   rawebscore.API_ROOT + "/notify",
	F: func(ttp *core.TTP) goutils.EchoRouteFunc {
		return func(c echo.Context) error {
			err := authenticateAdmin(ttp, c)
			if err != nil {
				return c.String(http.StatusUnauthorized, "token is invalid")
			}

			var data struct {
				Domain  string `json:"domain"`
				Message string `json:"message"`
			}

			err = c.Bind(&data)
			if err != nil {
				return err
			}

			serv, err := ttp.DB.Client.TAServer.
				Query().
				Where(taserver.DomainEQ(data.Domain)).
				Order(ent.Desc(taserver.FieldID)).
				First(*ttp.DB.Ctx)

			if err != nil {
				return err
			}

			err = ttp.Notify.Notify([]byte(data.Message), serv.Domain, ttp)
			if err != nil {
				return err
			}

			return c.JSON(http.StatusOK, "ok")
		}
	},
}
