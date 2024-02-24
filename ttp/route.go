package ttp

import (
	"fmt"
	"net/http"

	"github.com/akakou/ra_webs/ttp/api"
	"github.com/akakou/ra_webs/ttp/core"
	"github.com/akakou/ra_webs/ttp/ct"
	"github.com/akakou/ra_webs/ttp/web"
	"github.com/labstack/echo/v4"
)

func Route(e *echo.Echo, auditor *core.TTP) {
	e.GET("/", func(c echo.Context) error {
		r := fmt.Sprintf("%v", e.Routers())
		return c.String(http.StatusOK, r)
	})

	api.Route(e, auditor)
	ct.Route(e, auditor)
	web.Route(e, auditor)
}
