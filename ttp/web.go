package ttp

import (
	"net/http"

	"github.com/akakou/ra_webs/ttp/ent/taserver"
	"github.com/labstack/echo/v4"
)

var redirectWebPage = echoRoute{
	method: GET,
	path:   "/redirect",
	f: func(auditor *Auditor) echoRouteFunc {
		return func(c echo.Context) error {
			back := c.Request().Header.Get("Referer")

			server, err := auditor.db.Client.TAServer.Query().Where(taserver.DomainEQ(back)).WithTa().First(*auditor.db.Ctx)
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
