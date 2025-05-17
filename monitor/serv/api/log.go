package api

import (
	"fmt"
	"net/http"

	goutils "github.com/akakou/go-utils"
	"github.com/akakou/ra-webs/monitor/ent"
	"github.com/akakou/ra-webs/monitor/ent/ctlog"
	"github.com/akakou/ra-webs/monitor/serv"
	"github.com/labstack/echo/v4"
)

var GetLogs = goutils.EchoRoute[serv.MonitorServer]{
	Method: goutils.GET,
	Path:   "/ta",
	F: func(server *serv.MonitorServer) goutils.EchoRouteFunc {
		return func(c echo.Context) error {
			servs, err := server.Monitor.DB.Client.CTLog.
				Query().
				QueryTa().
				WithCtLog().
				WithAtLog().
				Order(ent.Desc(ctlog.FieldID)).
				All(*server.Monitor.DB.Ctx)

			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return c.JSON(http.StatusInternalServerError, "Internal Server Error")
			}

			return c.JSON(http.StatusOK, servs)
		}
	},
}
