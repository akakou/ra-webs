package ttp

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

var redirectWebPage = echoRoute{
	path: "/redirect",
	f: func(auditor *Auditor) echoRouteFunc {
		return func(c echo.Context) error {
			back := c.Request().Header.Get("Referer")

			return c.Render(http.StatusOK, "redirect", back)
		}
	},
}
