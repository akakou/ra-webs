package api

import (
	"fmt"
	"net/http"

	ttpcore "github.com/akakou/ra_webs/ttp/core"
	"github.com/labstack/echo/v4"
)

func Route(e *echo.Echo, ttp *ttpcore.TTP) {
	e.GET("/", func(c echo.Context) error {
		r := fmt.Sprintf("%v", e.Routers())
		return c.String(http.StatusOK, r)
	})

	PostCodeApi.Set(e, ttp)
	PostServerApi.Set(e, ttp)

	GetCodeApi.Set(e, ttp)
	GetServerApi.Set(e, ttp)

	PostServiceByAdmin.Set(e, ttp)
}
