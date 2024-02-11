package ttp

import (
	"crypto/x509/pkix"
	"net/http"

	"strconv"

	"github.com/akakou/ra_webs/core"
	"github.com/akakou/ra_webs/ttp/ent/ta"
	simplecertify "github.com/akakou/simple-certify"
	"github.com/labstack/echo/v4"
)

var registerTAApi = echoRoute{
	method: POST,
	path:   "/ta",
	f: func(auditor *Auditor) echoRouteFunc {
		return func(c echo.Context) error {
			reqTAInfo := new(struct {
				IP            string `json:"ip"`
				Domain        string `json:"domain"`
				GitRepository string `json:"git"`
			})

			if c.Bind(reqTAInfo) != nil {
				return c.String(http.StatusBadRequest, "bad attestation")
			}

			ta := auditor.db.Client.TA.
				Create().
				SetDomain(reqTAInfo.Domain).
				SetGit(reqTAInfo.GitRepository).
				SetIP(reqTAInfo.IP)

			t, err := ta.Save(*auditor.db.Ctx)

			if err != nil {
				c.Error(err)
			}

			err = auditor.ct.Subscribe(reqTAInfo.Domain)
			if err != nil {
				c.Error(err)
			}

			return c.String(http.StatusOK, strconv.Itoa(t.ID))
		}
	},
}

var updateTAApi = echoRoute{
	method: POST,
	path:   "/ta/:id/update",
	f: func(auditor *Auditor) echoRouteFunc {
		return func(c echo.Context) error {
			idParam := c.Param("id")

			id, err := strconv.Atoi(idParam)
			if err != nil {
				return c.String(http.StatusBadRequest, "bad id")
			}

			publicKey := new(struct {
				PublicKey []byte `json:"public_key"`
			})

			if c.Bind(publicKey) != nil {
				return c.String(http.StatusBadRequest, "bad public key")
			}

			ta, err := auditor.db.Client.TA.
				Query().Where(ta.IDEQ(id)).First(*auditor.db.Ctx)

			if err != nil {
				c.Error(err)
			}

			commitId, uniqueId := compile(ta)

			taCode := auditor.db.Client.TACode.
				Create().
				AddTa(ta).
				SetCommitID(commitId).
				SetUniqueID(uniqueId).
				SetPublicKey(publicKey.PublicKey)

			_, err = taCode.Save(*auditor.db.Ctx)
			if err != nil {
				c.Error(err)
			}

			templ := simplecertify.ServerTemplate()
			templ.PublicKey = publicKey.PublicKey
			templ.Subject = pkix.Name{
				Country:      []string{"Japan"},
				Organization: []string{"ra-webs"},
				Locality:     []string{"Kanagawa"},
				Province:     []string{"Yokohama"},
				CommonName:   ta.Domain,
			}
			templ.Issuer = auditor.ca.Certificate.Subject
			templ.Extensions = []pkix.Extension{
				{
					Id:    core.X509_EXTENSION_LABEL,
					Value: uniqueId,
				},
			}

			cert, err := auditor.ca.Certify(&templ)
			if err != nil {
				c.Error(err)
			}

			return c.Blob(http.StatusOK, "application/x-x509-cert", cert.Raw)
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
