package web

import (
	"net/http"

	goutils "github.com/akakou/go-utils"
	"github.com/akakou/ra_webs/ttp/core"
	"github.com/akakou/ra_webs/ttp/ent/taserver"
	"github.com/labstack/echo/v4"
)

var redirectWebPage = goutils.EchoRoute[core.TTP]{
	Method: goutils.GET,
	Path:   "/redirect",
	F: func(ttp *core.TTP) goutils.EchoRouteFunc {
		return func(c echo.Context) error {
			back := c.Request().Header.Get("Referer")

			server, err := ttp.DB.Client.TAServer.Query().Where(taserver.DomainEQ(back)).First(*ttp.DB.Ctx)

			if err != nil {
				c.Error(err)
				return err
			}

			ta, err := server.QueryTa().WithCtAudit().First(*ttp.DB.Ctx)
			if err != nil {
				c.Error(err)
				return err
			}

			if !ta.IsValid || !ta.Edges.CtAudit.IsValid {
				return c.String(http.StatusUnauthorized, "server is not valid")
			}

			return c.Render(http.StatusOK, "redirect", back)
		}
	},
}

func Route(e *echo.Echo, ttp *core.TTP) {
	redirectWebPage.Set(e, ttp)
}
