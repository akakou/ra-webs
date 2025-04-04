package api

import (
	"net/http"
	"strconv"

	goutils "github.com/akakou/go-utils"
	"github.com/akakou/ra-webs/log/core"
	"github.com/labstack/echo/v4"
)

var GetApi = goutils.EchoRoute[core.Log]{
	Method: goutils.GET,
	Path:   "/ta",
	F: func(log *core.Log) goutils.EchoRouteFunc {
		return func(c echo.Context) error {
			start := 0
			end := 0

			startParam := c.QueryParam("start")
			if startParam != "" {
				start, _ = strconv.Atoi(startParam)
				start = start - 1
			}

			endParam := c.QueryParam("end")
			if endParam != "" {
				end, _ = strconv.Atoi(endParam)
			}

			limit := end - start
			if limit <= 0 || limit > 100 {
				limit = 100
			}

			ta, err := log.DB.Client.TA.Query().Offset(start).Limit(limit).All(*log.DB.Ctx)

			if err != nil {
				return c.String(http.StatusBadRequest, "invalid request")
			}

			return c.JSON(http.StatusOK, ta)
		}
	},
}
