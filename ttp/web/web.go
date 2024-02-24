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
	F: func(auditor *core.TTP) goutils.EchoRouteFunc {
		return func(c echo.Context) error {
			back := c.Request().Header.Get("Referer")

			server, err := auditor.DB.Client.TAServer.Query().Where(taserver.DomainEQ(back)).WithTa().First(*auditor.DB.Ctx)
			if err != nil {
				c.Error(err)
				return err
			}

			if server.Edges.Ta == nil {
				return c.String(http.StatusUnauthorized, "ta is not found")
			}

			if !server.Edges.Ta.IsValid {
				return c.String(http.StatusUnauthorized, "server is not valid")
			}

			return c.Render(http.StatusOK, "redirect", back)
		}
	},
}

func Route(e *echo.Echo, auditor *core.TTP) {
	redirectWebPage.Set(e, auditor)
}
