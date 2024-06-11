package ct

import (
	"crypto/x509"
	"fmt"
	"net/http"

	goutils "github.com/akakou/go-utils"
	metact "github.com/akakou/meta-ct"
	rawebscore "github.com/akakou/ra_webs/core"
	"github.com/akakou/ra_webs/ttp/audit"
	"github.com/akakou/ra_webs/ttp/core"
	"github.com/labstack/echo/v4"
)

type MetaCT struct {
	Core metact.MetaCT
}

func NewMetaCT(metaAppId, metaAccessToken string) *MetaCT {
	mct := metact.NewCT(metaAppId, metaAccessToken)

	ct := &MetaCT{
		Core: *mct,
	}

	return ct
}

func (ct *MetaCT) Setup(e *echo.Echo, ttp *core.TTP) error {
	getWebhookApi, postWebhookApi := webhookApis()
	getWebhookApi.Set(e, ttp)
	postWebhookApi.Set(e, ttp)

	return nil

}

func (ct *MetaCT) Run(ttp *core.TTP) {
}

func (ct *MetaCT) Subscribe(domain string, ttp *core.TTP) error {
	return ct.Core.Subscribe(domain)
}

func metaCertsToCerts(cs []metact.MetaCert) ([]x509.Certificate, error) {
	certs := []x509.Certificate{}

	for _, c := range cs {
		cert, err := c.Certificate()

		if err != nil {
			return []x509.Certificate{}, err
		}

		certs = append(certs, *cert)
	}

	return certs, nil
}

func webhookApis() (goutils.EchoRoute[core.TTP], goutils.EchoRoute[core.TTP]) {
	hex, err := goutils.RandomHex(core.RANDOM_SIZE)
	if err != nil {
		panic(err)
	}
	path := rawebscore.API_ROOT + "/webhook/" + hex
	fmt.Printf("webhook path: %s\n", path)

	get := goutils.EchoRoute[core.TTP]{
		Method: goutils.GET,
		Path:   path,
		F: func(ttp *core.TTP) goutils.EchoRouteFunc {
			return func(c echo.Context) error {
				return metact.ChallengeAPIFlow(c)
			}
		},
	}

	post := goutils.EchoRoute[core.TTP]{
		Method: goutils.POST,
		Path:   path,
		F: func(ttp *core.TTP) goutils.EchoRouteFunc {
			return func(c echo.Context) error {
				m := ttp.CT.(*MetaCT)

				cs, err := m.Core.WebHookCertificates(c)
				if err != nil {
					return err
				}

				certs, err := metaCertsToCerts(cs)
				if err != nil {
					return err
				}

				err = audit.AuditAll(ttp, certs)
				if err != nil {
					return err
				}

				return c.String(http.StatusOK, "ok")
			}
		},
	}

	return get, post
}
