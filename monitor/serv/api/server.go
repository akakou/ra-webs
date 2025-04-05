package api

import (
	"net/http"

	goutils "github.com/akakou/go-utils"
	"github.com/akakou/ra-webs/monitor/ent"
	"github.com/akakou/ra-webs/monitor/ent/taserver"
	"github.com/akakou/ra-webs/monitor/serv"
	"github.com/labstack/echo/v4"
)

var GetServerApi = goutils.EchoRoute[serv.MonitorServer]{
	Method: goutils.GET,
	Path:   "/ta",
	F: func(server *serv.MonitorServer) goutils.EchoRouteFunc {
		return func(c echo.Context) error {
			code, err := server.Monitor.DB.Client.TAServer.Query().All(*server.Monitor.DB.Ctx)
			if err != nil {
				return err
			}

			return c.JSON(http.StatusOK, code)
		}
	},
}

var GetServerFromDomainApi = goutils.EchoRoute[serv.MonitorServer]{
	Method: goutils.GET,
	Path:   "/ta/:domain",
	F: func(server *serv.MonitorServer) goutils.EchoRouteFunc {
		return func(c echo.Context) error {
			type Res struct {
				TA      []*ent.TAServer `json:"ta"`
				IsValid bool            `json:"is_valid"`
				Message string          `json:"message"`
			}

			res := Res{}

			handleError := func(err error, r *Res) error {
				r.IsValid = false
				r.Message = err.Error()
				return c.JSON(http.StatusOK, res)
			}

			// fmt.Printf("domain: %v\n", domain)
			servs, err := server.Monitor.DB.Client.TAServer.
				Query().
				WithViolation().
				WithCode().
				WithService().
				Order(ent.Desc(taserver.FieldID)).
				All(*server.Monitor.DB.Ctx)

			res.TA = servs

			if err != nil {
				return handleError(err, &res)
			}

			isValid1, err := checkViolationLogs(servs)
			if err != nil {
				return handleError(err, &res)
			}

			isValid2, err := checkTAValidity(servs[0], server.Monitor)
			if err != nil {
				return handleError(err, &res)
			}

			res.IsValid = isValid1 && isValid2

			return c.JSON(http.StatusOK, res)
		}
	},
}
