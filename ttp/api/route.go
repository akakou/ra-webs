package api

import (
	"fmt"
	"net/http"

	ttpcore "github.com/akakou/ra_webs/ttp/core"
	"github.com/labstack/echo/v4"
)

const API_ROOT = "/api/"

func Route(e *echo.Echo, ttp *ttpcore.TTP) {
	e.GET("/", func(c echo.Context) error {
		r := fmt.Sprintf("%v", e.Routers())
		return c.String(http.StatusOK, r)
	})

	g := e.Group(API_ROOT)
	RegisterApi.Set(g, ttp)

	GetCodeApi.Set(g, ttp)
	GetServerApi.Set(g, ttp)

	PostServiceByAdmin.Set(g, ttp)

	GetServerFromDomainApi.Set(g, ttp)

	WebhookApi().Set(g, ttp)
}
