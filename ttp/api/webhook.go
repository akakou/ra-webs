package api

import (
	"fmt"
	"net/http"

	goutils "github.com/akakou/go-utils"
	"github.com/akakou/ra_webs/ttp/core"
	"github.com/akakou/ra_webs/ttp/ct"
	"github.com/labstack/echo/v4"
)

func WebhookApi() goutils.EchoRoute[core.TTP] {
	hex, err := goutils.RandomHex(core.RANDOM_SIZE)
	if err != nil {
		panic(err)
	}
	path := API_ROOT + "/webhook/" + hex
	fmt.Printf("webhook path: %s\n", path)

	return goutils.EchoRoute[core.TTP]{
		Method: goutils.POST,
		Path:   path,
		F: func(ttp *core.TTP) goutils.EchoRouteFunc {
			return func(c echo.Context) error {
				cs, err := ttp.CT.WebHookCertificates(c)
				if err != nil {
					return err
				}

				certs, err := ct.MetaCertsToCerts(cs)
				if err != nil {
					return err
				}

				err = ct.AuditAll(ttp, certs)
				if err != nil {
					return err
				}

				return c.String(http.StatusOK, "ok")
			}
		},
	}
}
