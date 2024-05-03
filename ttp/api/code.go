package api

import (
	"net/http"

	goutils "github.com/akakou/go-utils"
	ttpcore "github.com/akakou/ra_webs/ttp/core"
	"github.com/akakou/ra_webs/ttp/ent/tacode"
	"github.com/labstack/echo/v4"
)

var BUILD_DOCKER_PATH = "./builder"

var GetCodeApi = goutils.EchoRoute[ttpcore.TTP]{
	Method: goutils.GET,
	Path:   "/code",
	F: func(ttp *ttpcore.TTP) goutils.EchoRouteFunc {
		return func(c echo.Context) error {
			code, err := ttp.DB.Client.TACode.Query().
				Where(tacode.IsActive(true)).
				All(*ttp.DB.Ctx)
			if err != nil {
				return err
			}

			return c.JSON(http.StatusOK, code)
		}
	},
}
