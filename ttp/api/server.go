package api

import (
	"net/http"
	"strconv"

	goutils "github.com/akakou/go-utils"
	ttpcore "github.com/akakou/ra_webs/ttp/core"
	"github.com/akakou/ra_webs/ttp/ent/taserver"
	"github.com/labstack/echo/v4"
)

var postServerApi = goutils.EchoRoute[ttpcore.TTP]{
	Method: goutils.POST,
	Path:   "/server",
	F: func(ttp *ttpcore.TTP) goutils.EchoRouteFunc {
		return func(c echo.Context) error {
			service, err := authenticateService(ttp, c)
			if err != nil {
				return c.String(http.StatusUnauthorized, "token is invalid")
			}

			req := new(struct {
				IP     string `json:"ip"`
				Domain string `json:"domain"`
			})

			taServerCreate := ttp.DB.Client.TAServer.
				Create().
				SetIP(req.IP).
				SetDomain(req.Domain).
				SetService(service)

			taServer, err := taServerCreate.Save(*ttp.DB.Ctx)
			if err != nil {
				c.Error(err)
			}

			return c.String(http.StatusOK, strconv.Itoa(taServer.ID))
		}
	},
}

var postActivateServerApi = goutils.EchoRoute[ttpcore.TTP]{
	Method: goutils.POST,
	Path:   "/server/:id/activate",
	F: func(ttp *ttpcore.TTP) goutils.EchoRouteFunc {
		return func(c echo.Context) error {
			paramId := c.Param("id")
			codeId, err := strconv.Atoi(paramId)
			if err != nil {
				c.Error(err)
				return err
			}

			err = authenticateAdmin(ttp, c)
			if err != nil {
				return c.String(http.StatusUnauthorized, "token is invalid")
			}

			server, err := ttp.DB.Client.TAServer.Get(*ttp.DB.Ctx, codeId)
			if err != nil {
				c.Error(err)
				return err
			}

			_, err = server.Update().SetHasActivated(true).Save(*ttp.DB.Ctx)
			if err != nil {
				c.Error(err)
				return err
			}

			return c.String(http.StatusOK, strconv.Itoa(server.ID))
		}
	},
}

var getServerApi = goutils.EchoRoute[ttpcore.TTP]{
	Method: goutils.GET,
	Path:   "/server",
	F: func(ttp *ttpcore.TTP) goutils.EchoRouteFunc {
		return func(c echo.Context) error {
			activate := c.QueryParam("activate") != "false"

			code, err := ttp.DB.Client.TAServer.Query().Where(taserver.HasActivated(activate)).All(*ttp.DB.Ctx)
			if err != nil {
				c.Error(err)
				return err
			}

			return c.JSON(http.StatusOK, code)
		}
	},
}
