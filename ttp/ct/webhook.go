package ct

import (
	"fmt"
	"net/http"

	goutils "github.com/akakou/go-utils"
	"github.com/akakou/ra_webs/ttp/core"
	"github.com/labstack/echo/v4"
)

func webhook() goutils.EchoRoute[core.TTP] {
	hex, err := goutils.RandomHex(core.RANDOM_SIZE)
	if err != nil {
		panic(err)
	}
	path := "/webhook/" + hex
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

				certs, err := metaCertsToCerts(cs)
				if err != nil {
					return err
				}

				err = AuditAll(ttp, certs)
				if err != nil {
					return err
				}

				return c.String(http.StatusOK, "ok")
			}
		},
	}
}
