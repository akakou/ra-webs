package api

import (
	"net/http"

	goutils "github.com/akakou/go-utils"
	"github.com/akakou/ra-webs/monitor/serv"
	"github.com/labstack/echo/v4"
)

var PostNotifierApi = goutils.EchoRoute[serv.MonitorServer]{
	Method: goutils.POST,
	Path:   "/notify",
	F: func(server *serv.MonitorServer) goutils.EchoRouteFunc {
		return func(c echo.Context) error {
			err := authenticateAdmin(server, c)
			if err != nil {
				return c.String(http.StatusUnauthorized, "token is invalid")
			}

			var data struct {
				Domain  string `json:"domain"`
				Message string `json:"message"`
			}

			err = c.Bind(&data)
			if err != nil {
				return err
			}

			err = server.Monitor.Notifier.Notify([]byte(data.Message), data.Domain, server.Monitor)
			if err != nil {
				return err
			}

			return c.JSON(http.StatusOK, "ok")
		}
	},
}
