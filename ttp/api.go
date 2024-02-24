package ttp

import (
	"crypto/x509/pkix"
	"fmt"
	"net/http"

	"strconv"

	"github.com/akakou/ra_webs/core"
	simplecertify "github.com/akakou/simple-certify"
	"github.com/labstack/echo/v4"
)

var postCodeApi = echoRoute{
	method: POST,
	path:   "/code",
	f: func(auditor *Auditor) echoRouteFunc {
		return func(c echo.Context) error {
			req := new(struct {
				Repository string `json:"repository"`
				CommitId   string `json:"commit_id"`
				UniqueID   []byte `json:"unique_id"`
			})

			if c.Bind(req) != nil {
				return c.String(http.StatusBadRequest, "bad attestation")
			}

			codeCreate := auditor.db.Client.TACode.
				Create().
				SetRepository(req.Repository).
				SetCommitID(req.CommitId).
				SetUniqueID(req.UniqueID)

			code, err := codeCreate.Save(*auditor.db.Ctx)

			if err != nil {
				c.Error(err)
			}

			return c.String(http.StatusOK, strconv.Itoa(code.ID))
		}
	},
}

var postServerApi = echoRoute{
	method: POST,
	path:   "/server",
	f: func(auditor *Auditor) echoRouteFunc {
		return func(c echo.Context) error {
			serviceId := "000"

			req := new(struct {
				IP     string `json:"ip"`
				Domain string `json:"domain"`
			})

			taServerCreate := auditor.db.Client.TAServer.
				Create().
				SetIP(req.IP).
				SetDomain(req.Domain).
				SetServiceID(serviceId)

			taServer, err := taServerCreate.Save(*auditor.db.Ctx)
			if err != nil {
				c.Error(err)
			}

			fmt.Printf("id 1: %v", taServer.ID)

			return c.String(200, strconv.Itoa(taServer.ID))
		}
	},
}

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

			code, err := auditor.db.Client.TACode.Get(*auditor.db.Ctx, req.CodeId)
			if err != nil {
				c.Error(err)
				return err
			}

			serv, err := auditor.db.Client.TAServer.Get(*auditor.db.Ctx, req.ServerId)

			if err != nil {
				c.Error(err)
				return err
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
