package ttp

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (auditServ *AuditServer) webhook() echoRoute {
	return func(c echo.Context) error {
		certs, err := auditServ.auditor.ct.WebHookCertificates(c)
		if err != nil {
			c.Error(err)
		}

		err = auditServ.auditor.AuditAll(certs)
		if err != nil {
			c.Error(err)
		}

		return c.String(http.StatusOK, "ok")
	}
}
