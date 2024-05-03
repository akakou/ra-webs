package api

import (
	"errors"
	"fmt"
	"net/http"

	goutils "github.com/akakou/go-utils"
	ttpcore "github.com/akakou/ra_webs/ttp/core"
	"github.com/akakou/ra_webs/ttp/ent/taserver"
	"github.com/labstack/echo/v4"
)

var GetServerApi = goutils.EchoRoute[ttpcore.TTP]{
	Method: goutils.GET,
	Path:   "/server",
	F: func(ttp *ttpcore.TTP) goutils.EchoRouteFunc {
		return func(c echo.Context) error {
			code, err := ttp.DB.Client.TAServer.Query().All(*ttp.DB.Ctx)
			if err != nil {
				return err
			}

			return c.JSON(http.StatusOK, code)
		}
	},
}

var GetServerFromDomainApi = goutils.EchoRoute[ttpcore.TTP]{
	Method: goutils.GET,
	Path:   "/server/:domain",
	F: func(ttp *ttpcore.TTP) goutils.EchoRouteFunc {
		return func(c echo.Context) error {
			domain := c.Param("domain")

			fmt.Printf("domain: %v\n", domain)
			servs, err := ttp.DB.Client.TAServer.
				Query().
				Where(taserver.Domain(domain)).
				Where(taserver.HasActivated(true)).
				WithCode().
				WithViolation().
				All(*ttp.DB.Ctx)

			if err != nil {
				return errors.New("server is not found")
			}

			return c.JSON(http.StatusOK, servs)
		}
	},
}
