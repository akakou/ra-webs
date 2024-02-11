package ttp

import (
	"fmt"
	"net/http"

	goutils "github.com/akakou/go-utils"
	"github.com/labstack/echo/v4"
)

func webhook() echoRoute {
	hex, err := goutils.RandomHex(RANDOM_SIZE)
	if err != nil {
		panic(err)
	}
	path := "/webhook/" + hex
	fmt.Printf("webhook path: %s\n", path)

	return echoRoute{
		method: POST,
		path:   path,
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
