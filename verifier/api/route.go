package api

import (
	"fmt"
	"net/http"

	verifiercore "github.com/akakou/ra_webs/verifier/core"
	"github.com/labstack/echo/v4"
)

func Route(e *echo.Echo, verifier *verifiercore.Verifier) {
	e.GET("/", func(c echo.Context) error {
		r := fmt.Sprintf("%v", e.Routers())
		return c.String(http.StatusOK, r)
	})

	RegisterApi.Set(e, verifier)
	GetServerApi.Set(e, verifier)

	PostServiceByAdmin.Set(e, verifier)

	GetServerFromDomainApi.Set(e, verifier)
	PostNotifierApi.Set(e, verifier)
}
