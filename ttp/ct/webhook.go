package ct

import (
	"fmt"
	"net/http"

	goutils "github.com/akakou/go-utils"
	"github.com/akakou/ra_webs/core"
	ttpcore "github.com/akakou/ra_webs/ttp/core"
	"github.com/labstack/echo/v4"
)

func webhook() goutils.EchoRoute[ttpcore.TTP] {
	hex, err := goutils.RandomHex(core.RANDOM_SIZE)
	if err != nil {
		panic(err)
	}
	path := "/webhook/" + hex
	fmt.Printf("webhook path: %s\n", path)

	return goutils.EchoRoute[ttpcore.TTP]{
		Method: goutils.POST,
		Path:   path,
		F: func(ttp *ttpcore.TTP) goutils.EchoRouteFunc {
			return func(c echo.Context) error {
				certs, err := ttp.CT.WebHookCertificates(c)
				if err != nil {
					c.Error(err)
				}

				err = AuditAll(ttp, certs)
				if err != nil {
					c.Error(err)
				}

				return c.String(http.StatusOK, "ok")
			}
		},
	}
}
