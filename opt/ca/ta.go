package ca

import (
	"net/http"
	"strconv"

	goutils "github.com/akakou/go-utils"
	"github.com/labstack/echo"
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

			ta, err := ttp.DB.Client.TA.Create().SetPublicKey(req.PublicKey).SetCode(code).SetIsValid(false).SetServer(serv).Save(*ttp.DB.Ctx)
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
	Path:   "/ta/:id/start",
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

			req := new(struct {
				Quote string `json:"quote"`
			})

			err = c.Bind(req)
			if err != nil {
				c.Error(err)
				return err
			}

			serv, err := ta.QueryServer().WithService().Only(*ttp.DB.Ctx)
			if err != nil {
				c.Error(err)
				return err
			}

			if serv.Edges.Service.ID != service.ID {
				return c.String(http.StatusUnauthorized, "ta is not authorized")
			}

			code, err := ta.QueryCode().Only(*ttp.DB.Ctx)
			if err != nil {
				c.Error(err)
				return err
			}

			_, err = core.VerifyByAzure(req.Quote, ta.PublicKey, code.UniqueID)
			if err != nil {
				c.Error(err)
				return err
			}

			cert, err := issueCertificate(serv.Domain, code.UniqueID, ttp.CA)
			if err != nil {
				c.Error(err)
				return err
			}

			ta.Update().SetIsValid(true).Save(*ttp.DB.Ctx)

			resp := map[string]interface{}{
				"cert":      cert,
				"unique_id": code.UniqueID,
			}

			return c.JSON(http.StatusOK, resp)
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
