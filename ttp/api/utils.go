package api

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"net/http"

	"github.com/akakou/ra_webs/core"
	ttpcore "github.com/akakou/ra_webs/ttp/core"
	"github.com/akakou/ra_webs/ttp/ent"
	"github.com/akakou/ra_webs/ttp/ent/service"
	"github.com/akakou/ra_webs/ttp/ent/ta"
	simplecertify "github.com/akakou/simple-certify"
	"github.com/labstack/echo/v4"
)

func authenticateService(ttp *ttpcore.TTP, c echo.Context) (*ent.Service, error) {
	authorization := c.Request().Header["Authorization"][0]
	token := authorization[len("Bearer "):]

	service, err := ttp.DB.Client.Service.Query().Where(service.TokenEQ(token)).First(*ttp.DB.Ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate service: %w", err)
	}

	return service, nil
}

func authenticateAdmin(ttp *ttpcore.TTP, c echo.Context) error {
	authorization := c.Request().Header["Authorization"][0]
	token := authorization[len("Bearer "):]

	if token != ttp.AdminToken {
		return c.String(http.StatusUnauthorized, "token is invalid")
	}

	return nil
}

func issueCertificate(domain string, uniqueId []byte, ca *simplecertify.Certifier) (*x509.Certificate, error) {
	templ := simplecertify.ServerTemplate()
	templ.PublicKey = ta.PublicKey
	templ.Subject = pkix.Name{
		Country:      []string{"Japan"},
		Organization: []string{"ra-webs"},
		Locality:     []string{"Kanagawa"},
		Province:     []string{"Yokohama"},
		CommonName:   domain,
	}

	templ.Issuer = ca.Certificate.Subject
	templ.Extensions = []pkix.Extension{
		{
			Id:    core.X509_EXTENSION_LABEL,
			Value: []byte(uniqueId),
		},
	}

	cert, err := ca.Certify(&templ)
	if err != nil {
		return nil, fmt.Errorf("failed to issue certificate: %w", err)
	}

	return cert, nil

}
