package api

import (
	"crypto/x509/pkix"
	"fmt"
	"net/http"

	goutils "github.com/akakou/go-utils"
	"github.com/akakou/ra_webs/core"
	ttpcore "github.com/akakou/ra_webs/ttp/core"

	"github.com/akakou/ra_webs/ttp/ent/ta"
	"github.com/akakou/ra_webs/ttp/ent/taserver"
	simplecertify "github.com/akakou/simple-certify"
	"github.com/labstack/echo/v4"
)

func Route(e *echo.Echo, ttp *ttpcore.TTP) {
	e.GET("/", func(c echo.Context) error {
		r := fmt.Sprintf("%v", e.Routers())
		return c.String(http.StatusOK, r)
	})

	postCodeApi.Set(e, ttp)
	postServerApi.Set(e, ttp)
	postTAApi.Set(e, ttp)
	certApi.Set(e, ttp)

	getCodeApi.Set(e, ttp)
	getServerApi.Set(e, ttp)
	getTAApi.Set(e, ttp)

	postActivateServerApi.Set(e, ttp)
	postActivateCodeApi.Set(e, ttp)

	postServiceByAdmin.Set(e, ttp)
}

var postTAApi = goutils.EchoRoute[ttpcore.TTP]{
	Method: goutils.POST,
	Path:   "/ta",
	F: func(ttp *ttpcore.TTP) goutils.EchoRouteFunc {
		return func(c echo.Context) error {
			service, err := authenticateService(ttp, c)

			if err != nil {
				return c.String(http.StatusUnauthorized, "token is invalid")
			}

			req := new(struct {
				PublicKey []byte `json:"public_key"`
				CodeId    int    `json:"code_id"`
				ServerId  int    `json:"server_id"`
			})

			err = c.Bind(req)
			if err != nil {
				c.Error(err)
				return err
			}

			serv, err := ttp.DB.Client.TAServer.Query().WithService().Where(taserver.ID(req.ServerId)).Only(*ttp.DB.Ctx)

			if err != nil {
				c.Error(err)
				return err
			}

			if !serv.HasActivated {
				return c.String(http.StatusUnauthorized, "server is not activated")
			}

			if serv.Edges.Service.ID != service.ID {
				return c.String(http.StatusUnauthorized, "server is not authorized")
			}

			code, err := ttp.DB.Client.TACode.Get(*ttp.DB.Ctx, req.CodeId)
			if err != nil {
				c.Error(err)
				return err
			}

			if !code.HasActivated {
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

			templ.Issuer = ttp.CA.Certificate.Subject
			templ.Extensions = []pkix.Extension{
				{
					Id:    core.X509_EXTENSION_LABEL,
					Value: code.UniqueID,
				},
			}

			cert, err := ttp.CA.Certify(&templ)
			if err != nil {
				c.Error(err)
				return err
			}

			// err = ttp.ct.Subscribe(serv.Domain)
			// if err != nil {
			// 	c.Error(err)
			// 	return err
			// }

			return c.Blob(http.StatusOK, "application/x-x509-cert", cert.Raw)
		}
	},
}

var getTAApi = goutils.EchoRoute[ttpcore.TTP]{
	Method: goutils.GET,
	Path:   "/ta",
	F: func(ttp *ttpcore.TTP) goutils.EchoRouteFunc {
		return func(c echo.Context) error {
			valid := c.QueryParam("valid") != "false"

			ta, err := ttp.DB.Client.TA.Query().Where(ta.IsValid(valid)).All(*ttp.DB.Ctx)
			if err != nil {
				c.Error(err)
				return err
			}

			return c.JSON(http.StatusOK, ta)
		}
	},
}

var certApi = goutils.EchoRoute[ttpcore.TTP]{
	Method: goutils.GET,
	Path:   "/cert",
	F: func(ttp *ttpcore.TTP) goutils.EchoRouteFunc {
		return func(c echo.Context) error {
			cert := ttp.CA.Certificate.Raw
			return c.Blob(http.StatusOK, "application/x-x509-ca-cert", cert)
		}
	},
}
