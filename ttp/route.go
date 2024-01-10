package ttp

import (
	"net/http"

	"github.com/akakou/ra_webs/core"
	"github.com/labstack/echo/v4"
)

func Route(e *echo.Echo, db *ttpDB) {
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.POST("/provision", func(c echo.Context) error {
		provReq := new(core.TAInfo)

		if err := c.Bind(provReq); err != nil {
			return err
		}

		if verifyAttestation(provReq.Attestation) != nil {
			return c.String(http.StatusBadRequest, "bad attestation")
		}

		if db.store(provReq) != nil {
			return c.String(http.StatusInternalServerError, "internal error")
		}

		return c.String(http.StatusOK, "ok")
	})

	e.GET("/redirect", func(c echo.Context) error {
		back := c.Request().Header.Get("Referer")
		// back = "aaa"
		return c.Render(http.StatusOK, "redirect", back)
	})
}
