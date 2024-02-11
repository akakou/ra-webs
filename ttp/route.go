package ttp

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

var RANDOM_SIZE = 32

type echoRouteFunc = func(c echo.Context) error

type echoRoute struct {
	path string
	f    func(*Auditor) echoRouteFunc
}

func (er echoRoute) get(e *echo.Echo, auditor *Auditor) {
	e.GET(er.path, er.f(auditor))
}

func (er echoRoute) post(e *echo.Echo, auditor *Auditor) {
	e.POST(er.path, er.f(auditor))
}

func Route(e *echo.Echo, auditor *Auditor) {
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	registerTAApi.post(e, auditor)
	updateTAApi.post(e, auditor)
	certApi.get(e, auditor)

	webhook().post(e, auditor)

	redirectWebPage.get(e, auditor)
}
