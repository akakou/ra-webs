package ttp

import (
	"net/http"
	"strconv"

	"github.com/akakou/ra_webs/ttp/ent/taserver"
	"github.com/labstack/echo/v4"
)

var postServerApi = echoRoute{
	method: POST,
	path:   "/server",
	f: func(auditor *Auditor) echoRouteFunc {
		return func(c echo.Context) error {
			service, err := authenticateService(auditor.db, c)
			if err != nil {
				return c.String(http.StatusUnauthorized, "token is invalid")
			}

			req := new(struct {
				IP     string `json:"ip"`
				Domain string `json:"domain"`
			})

			taServerCreate := auditor.db.Client.TAServer.
				Create().
				SetIP(req.IP).
				SetDomain(req.Domain).
				SetService(service)

			taServer, err := taServerCreate.Save(*auditor.db.Ctx)
			if err != nil {
				c.Error(err)
			}

			return c.String(http.StatusOK, strconv.Itoa(taServer.ID))
		}
	},
}

var postActivateServerApi = echoRoute{
	method: POST,
	path:   "/server/:id/activate",
	f: func(auditor *Auditor) echoRouteFunc {
		return func(c echo.Context) error {
			paramId := c.Param("id")
			codeId, err := strconv.Atoi(paramId)
			if err != nil {
				c.Error(err)
				return err
			}

			err = authenticateAdmin(auditor, c)
			if err != nil {
				return c.String(http.StatusUnauthorized, "token is invalid")
			}

			server, err := auditor.db.Client.TAServer.Get(*auditor.db.Ctx, codeId)
			if err != nil {
				c.Error(err)
				return err
			}

			_, err = server.Update().SetHasActivated(true).Save(*auditor.db.Ctx)
			if err != nil {
				c.Error(err)
				return err
			}

			return c.String(http.StatusOK, strconv.Itoa(server.ID))
		}
	},
}

var getServerApi = echoRoute{
	method: GET,
	path:   "/server",
	f: func(auditor *Auditor) echoRouteFunc {
		return func(c echo.Context) error {
			activate := c.QueryParam("activate") != "false"

			code, err := auditor.db.Client.TAServer.Query().Where(taserver.HasActivated(activate)).All(*auditor.db.Ctx)
			if err != nil {
				c.Error(err)
				return err
			}

			return c.JSON(http.StatusOK, code)
		}
	},
}
