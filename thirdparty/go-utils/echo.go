package goutils

import "github.com/labstack/echo/v4"

type EchoRouteFunc = func(c echo.Context) error

type EchoRoute[T any] struct {
	Method int
	Path   string
	F      func(*T) EchoRouteFunc
}

const (
	ANY = iota
	GET
	POST
)

func (er EchoRoute[T]) Set(e *echo.Group, param *T) {
	if er.Method == ANY {
		e.Any(er.Path, er.F(param))
	}
	if er.Method == GET {
		e.GET(er.Path, er.F(param))
	}
	if er.Method == POST {
		e.POST(er.Path, er.F(param))
	}
}
