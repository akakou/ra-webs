package ttp

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Route(e *echo.Echo, db *auditDB) {
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.POST("/register", func(c echo.Context) error {
		reqTAInfo := new(RegisterReqBody)

		if c.Bind(reqTAInfo) != nil {
			return c.String(http.StatusBadRequest, "bad attestation")
		}

		taInfo := db.client.TAInfo.
			Create().
			SetDomain(reqTAInfo.Domain).
			SetGitRepository(reqTAInfo.GitRepository)

		_, err := taInfo.Save(*db.ctx)
		if err != nil {
			return c.String(http.StatusInternalServerError, "internal error")
		}

		return c.String(http.StatusOK, "ok")
	})

	e.GET("/redirect", func(c echo.Context) error {
		back := c.Request().Header.Get("Referer")

		return c.Render(http.StatusOK, "redirect", back)
	})
}
