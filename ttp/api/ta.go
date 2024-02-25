package api

import (
	"crypto/x509/pkix"
	"net/http"
	"strconv"

	goutils "github.com/akakou/go-utils"
	"github.com/akakou/ra_webs/core"
	ttpcore "github.com/akakou/ra_webs/ttp/core"
	"github.com/akakou/ra_webs/ttp/ent/ta"
	"github.com/akakou/ra_webs/ttp/ent/taserver"
	simplecertify "github.com/akakou/simple-certify"
	"github.com/labstack/echo/v4"
)

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

			ta, err := ttp.DB.Client.TA.Create().SetPublicKey(req.PublicKey).SetCode(code).SetLastCt("").SetServer(serv).Save(*ttp.DB.Ctx)
			if err != nil {
				c.Error(err)
				return err
			}

			// err = ttp.ct.Subscribe(serv.Domain)
			// if err != nil {
			// 	c.Error(err)
			// 	return err
			// }

			return c.String(http.StatusOK, strconv.Itoa(ta.ID))
		}
	},
}

var getTACertApi = goutils.EchoRoute[ttpcore.TTP]{
	Method: goutils.GET,
	Path:   "/ta/:id/cert",
	F: func(ttp *ttpcore.TTP) goutils.EchoRouteFunc {
		return func(c echo.Context) error {
			paramId := c.Param("id")
			taId, err := strconv.Atoi(paramId)
			if err != nil {
				c.Error(err)
				return err
			}

			service, err := authenticateService(ttp, c)
			if err != nil {
				return c.String(http.StatusUnauthorized, "token is invalid")
			}

			ta, err := ttp.DB.Client.TA.Get(*ttp.DB.Ctx, taId)
			if err != nil {
				c.Error(err)
				return err
			}

			serv, err := ta.QueryServer().WithService().First(*ttp.DB.Ctx)
			if err != nil {
				c.Error(err)
				return err
			}

			if serv.Edges.Service.ID != service.ID {
				return c.String(http.StatusUnauthorized, "ta is not authorized")
			}

			code, err := ta.QueryCode().First(*ttp.DB.Ctx)
			if err != nil {
				c.Error(err)
				return err
			}

			templ := simplecertify.ServerTemplate()
			templ.PublicKey = ta.PublicKey
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
