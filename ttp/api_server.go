package ttp

import (
	"net/http"
	"strconv"

	goutils "github.com/akakou/go-utils"
	"github.com/akakou/ra_webs/ttp/ent/taserver"
	"github.com/labstack/echo/v4"
)

var postServerApi = echoRoute{
	method: POST,
	path:   "/server",
	f: func(auditor *Auditor) echoRouteFunc {
		return func(c echo.Context) error {
			serviceId := "000"

			req := new(struct {
				IP     string `json:"ip"`
				Domain string `json:"domain"`
			})

			token, err := goutils.RandomHex(32)
			if err != nil {
				c.Error(err)
				return err
			}

			taServerCreate := auditor.db.Client.TAServer.
				Create().
				SetIP(req.IP).
				SetDomain(req.Domain).
				SetServiceID(serviceId).
				SetToken(token)

			taServer, err := taServerCreate.Save(*auditor.db.Ctx)
			if err != nil {
				c.Error(err)
			}

			return c.JSON(200,
				map[string]interface{}{
					"server_id": taServer.ID,
					"token":     taServer.Token,
				})
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

			authorization := c.Request().Header["Authorization"][0]
			token := authorization[len("Bearer "):]

			if token != auditor.adminToken {
				return c.String(http.StatusUnauthorized, "token is invalid")
			}

			server, err := auditor.db.Client.TAServer.Get(*auditor.db.Ctx, codeId)
			if err != nil {
				c.Error(err)
				return err
			}

			_, err = server.Update().SetActivate(true).Save(*auditor.db.Ctx)
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

			code, err := auditor.db.Client.TAServer.Query().Where(taserver.Activate(activate)).All(*auditor.db.Ctx)
			if err != nil {
				c.Error(err)
				return err
			}

			return c.JSON(http.StatusOK, code)
		}
	},
}
