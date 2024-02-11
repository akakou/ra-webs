package ttp

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

var RANDOM_SIZE = 32

type echoRouteFunc = func(c echo.Context) error

type echoRoute struct {
	method int
	path   string
	f      func(*Auditor) echoRouteFunc
}

const (
	ANY = iota
	GET
	POST
)

func (er echoRoute) set(e *echo.Echo, auditor *Auditor) {
	if er.method == ANY {
		e.Any(er.path, er.f(auditor))
	}
	if er.method == GET {
		e.GET(er.path, er.f(auditor))
	}
	if er.method == POST {
		e.POST(er.path, er.f(auditor))
	}
}

func Route(e *echo.Echo, auditor *Auditor) {
	e.GET("/", func(c echo.Context) error {
		r := fmt.Sprintf("%v", e.Routers())
		return c.String(http.StatusOK, r)
	})

	registerTAApi.set(e, auditor)
	updateTAApi.set(e, auditor)
	certApi.set(e, auditor)

	webhook().set(e, auditor)

	redirectWebPage.set(e, auditor)
}
