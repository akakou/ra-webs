package ttp

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

var RANDOM_SIZE = 32

type echoRoute = func(c echo.Context) error

type AuditServer struct {
	auditor *Auditor
}

func NewAuditServer(auditor *Auditor) *AuditServer {
	return &AuditServer{auditor}
}

func (auditServ *AuditServer) Route(e *echo.Echo) {
	webhookPath := "/webhook/" + randomHexString(RANDOM_SIZE)
	fmt.Printf("webhook path: %s\n", webhookPath)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.POST("/register", auditServ.register())

	e.POST("/compile", auditServ.compile())

	e.GET("/webhook", auditServ.webhook())

	e.GET("/redirect", func(c echo.Context) error {
		back := c.Request().Header.Get("Referer")

		return c.Render(http.StatusOK, "redirect", back)
	})
}
