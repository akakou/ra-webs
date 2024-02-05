package ttp

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

var RANDOM_SIZE = 32

func Route(e *echo.Echo, auditor *Auditor) {
	webhookPath := "/webhook/" + randomHexString(RANDOM_SIZE)
	fmt.Printf("webhook path: %s\n", webhookPath)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.POST("/register", func(c echo.Context) error {
		reqTAInfo := new(struct {
			Domain        string
			GitRepository string
		})

		if c.Bind(reqTAInfo) != nil {
			return c.String(http.StatusBadRequest, "bad attestation")
		}

		taInfo := auditor.db.client.TAInfo.
			Create().
			SetDomain(reqTAInfo.Domain).
			SetGitRepository(reqTAInfo.GitRepository)

		_, err := taInfo.Save(*auditor.db.ctx)
		if err != nil {
			c.Error(err)
		}

		return c.String(http.StatusOK, "ok")
	})

	e.GET(webhookPath, func(c echo.Context) error {
		certs, err := auditor.ct.WebHookCertificates(c)
		if err != nil {
			c.Error(err)
		}

		err = auditor.AuditAll(certs)
		if err != nil {
			c.Error(err)
		}

		return c.String(http.StatusOK, "ok")
	})

	e.GET("/redirect", func(c echo.Context) error {
		back := c.Request().Header.Get("Referer")

		return c.Render(http.StatusOK, "redirect", back)
	})
}
