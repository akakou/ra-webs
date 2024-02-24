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

	postCodeApi.Set(e, ttp)
	postServerApi.Set(e, ttp)
	postTAApi.Set(e, ttp)
	certApi.Set(e, ttp)

	getCodeApi.Set(e, ttp)
	getServerApi.Set(e, ttp)
	getTAApi.Set(e, ttp)

	postActivateServerApi.Set(e, ttp)
	postActivateCodeApi.Set(e, ttp)

	postServiceByAdmin.Set(e, ttp)
}
