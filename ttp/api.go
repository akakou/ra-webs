package ttp

import (
	"crypto/x509/pkix"
	"net/http"

	"github.com/akakou/ra_webs/core"
	"github.com/akakou/ra_webs/ttp/ent/ta"
	simplecertify "github.com/akakou/simple-certify"
	"github.com/labstack/echo/v4"
)

var postTAApi = echoRoute{
	method: POST,
	path:   "/ta",
	f: func(auditor *Auditor) echoRouteFunc {
		return func(c echo.Context) error {
			req := new(struct {
				PublicKey []byte `json:"public_key"`
				CodeId    int    `json:"code_id"`
				ServerId  int    `json:"server_id"`
			})

			err := c.Bind(req)
			if err != nil {
				c.Error(err)
				return err
			}

			serv, err := auditor.db.Client.TAServer.Get(*auditor.db.Ctx, req.ServerId)

			if err != nil {
				c.Error(err)
				return err
			}

			if !serv.Activate {
				return c.String(http.StatusUnauthorized, "server is not activated")
			}

			authorization := c.Request().Header["Authorization"][0]
			token := authorization[len("Bearer "):]

			if token != serv.Token {
				return c.String(http.StatusUnauthorized, "token is invalid")
			}

			code, err := auditor.db.Client.TACode.Get(*auditor.db.Ctx, req.CodeId)
			if err != nil {
				c.Error(err)
				return err
			}

			if !code.Activate {
				return c.String(http.StatusUnauthorized, "code is not activated")
			}

			templ := simplecertify.ServerTemplate()
			templ.PublicKey = req.PublicKey
			templ.Subject = pkix.Name{
				Country:      []string{"Japan"},
				Organization: []string{"ra-webs"},
				Locality:     []string{"Kanagawa"},
				Province:     []string{"Yokohama"},
				CommonName:   serv.Domain,
			}

			templ.Issuer = auditor.ca.Certificate.Subject
			templ.Extensions = []pkix.Extension{
				{
					Id:    core.X509_EXTENSION_LABEL,
					Value: code.UniqueID,
				},
			}

			cert, err := auditor.ca.Certify(&templ)
			if err != nil {
				c.Error(err)
				return err
			}

			// err = auditor.ct.Subscribe(serv.Domain)
			// if err != nil {
			// 	c.Error(err)
			// 	return err
			// }

			return c.Blob(http.StatusOK, "application/x-x509-cert", cert.Raw)
		}
	},
}

var getTAApi = echoRoute{
	method: GET,
	path:   "/ta",
	f: func(auditor *Auditor) echoRouteFunc {
		return func(c echo.Context) error {
			valid := c.QueryParam("valid") != "false"

			ta, err := auditor.db.Client.TA.Query().Where(ta.IsValid(valid)).All(*auditor.db.Ctx)
			if err != nil {
				c.Error(err)
				return err
			}

			return c.JSON(http.StatusOK, ta)
		}
	},
}

var certApi = echoRoute{
	method: GET,
	path:   "/cert",
	f: func(auditor *Auditor) echoRouteFunc {
		return func(c echo.Context) error {
			cert := auditor.ca.Certificate.Raw
			return c.Blob(http.StatusOK, "application/x-x509-ca-cert", cert)
		}
	},
}
