package ttp

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func webhook() echoRoute {
	path := "/webhook/" + randomHexString(RANDOM_SIZE)
	fmt.Printf("webhook path: %s\n", path)

	return echoRoute{
		path: path,
		f: func(auditor *Auditor) echoRouteFunc {
			return func(c echo.Context) error {
				certs, err := auditor.ct.WebHookCertificates(c)
				if err != nil {
					c.Error(err)
				}

				err = auditor.AuditAll(certs)
				if err != nil {
					c.Error(err)
				}

				return c.String(http.StatusOK, "ok")
			}
		},
	}
}
