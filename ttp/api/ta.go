package api

import (
	"errors"
	"fmt"
	"net/http"

	goutils "github.com/akakou/go-utils"
	ttpcore "github.com/akakou/ra_webs/ttp/core"
	"github.com/akakou/ra_webs/ttp/ent/ta"
	"github.com/akakou/ra_webs/ttp/ent/taserver"
	"github.com/labstack/echo/v4"
)

var GetTAFromDomainApi = goutils.EchoRoute[ttpcore.TTP]{
	Method: goutils.GET,
	Path:   "/ta/:domain",
	F: func(ttp *ttpcore.TTP) goutils.EchoRouteFunc {
		return func(c echo.Context) error {
			domain := c.Param("domain")

			fmt.Printf("domain: %v\n", domain)
			serv, err := ttp.DB.Client.TAServer.
				Query().
				Where(taserver.Domain(domain)).
				Only(*ttp.DB.Ctx)

			if err != nil {
				return errors.New("server is not found")
			}

			tas, err := ttp.DB.Client.TA.
				Query().
				Where(
					ta.HasServerWith(taserver.ID(
						serv.ID,
					)),
				).
				WithCode().
				All(*ttp.DB.Ctx)

			if err != nil {
				return errors.New("ta is not found")
			}

			return c.JSON(http.StatusOK, tas)
		}
	},
}
