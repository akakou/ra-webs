package api

import (
	"net/http"

	goutils "github.com/akakou/go-utils"
	rawebscore "github.com/akakou/ra_webs/core"
	"github.com/akakou/ra_webs/verifier/core"
	"github.com/akakou/ra_webs/verifier/ent"
	"github.com/akakou/ra_webs/verifier/ent/taserver"
	"github.com/labstack/echo/v4"
)

var PostNotifierApi = goutils.EchoRoute[core.Verifier]{
	Method: goutils.POST,
	Path:   rawebscore.API_ROOT + "/notifier",
	F: func(verifier *core.Verifier) goutils.EchoRouteFunc {
		return func(c echo.Context) error {
			err := authenticateAdmin(verifier, c)
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

			serv, err := verifier.DB.Client.TAServer.
				Query().
				Where(taserver.DomainEQ(data.Domain)).
				Order(ent.Desc(taserver.FieldID)).
				First(*verifier.DB.Ctx)

			if err != nil {
				return err
			}

			err = verifier.Notifier.Notifier([]byte(data.Message), serv.Domain, verifier)
			if err != nil {
				return err
			}

			return c.JSON(http.StatusOK, "ok")
		}
	},
}
