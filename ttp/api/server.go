package api

import (
	"fmt"
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
			fmt.Print("1\n")
			service, err := authenticateService(ttp, c)
			if err != nil {
				return c.String(http.StatusUnauthorized, "token is invalid")
			}

			req := new(struct {
				Domain string `json:"domain"`
				Nonce  string `json:"nonce"`
			})

			if c.Bind(req) != nil {
				return c.String(http.StatusBadRequest, "bad request")
			}

			err = authenticateDomain(req.Domain, service.Token, req.Nonce)
			if err != nil {
				return c.String(http.StatusUnauthorized, err.Error())
			}

			taServerCreate := ttp.DB.Client.TAServer.
				Create().
				SetDomain(req.Domain).
				SetService(service).
				SetIsActive(true)

			taServer, err := taServerCreate.Save(*ttp.DB.Ctx)
			if err != nil {
				return c.String(http.StatusInternalServerError, err.Error())
			}

			return c.String(http.StatusOK, strconv.Itoa(taServer.ID))
		}
	},
}

var getServerApi = goutils.EchoRoute[ttpcore.TTP]{
	Method: goutils.GET,
	Path:   "/server",
	F: func(ttp *ttpcore.TTP) goutils.EchoRouteFunc {
		return func(c echo.Context) error {
			code, err := ttp.DB.Client.TAServer.Query().Where(taserver.IsActive(true)).All(*ttp.DB.Ctx)
			if err != nil {
				return err
			}

			return c.JSON(http.StatusOK, code)
		}
	},
}
