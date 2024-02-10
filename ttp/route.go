package ttp

import (
	"fmt"
	"net/http"

	"github.com/akakou/ra_webs/ttp/ent/tainfo"
	"github.com/labstack/echo/v4"
)

var RANDOM_SIZE = 32

func Route(e *echo.Echo, auditor *Auditor) {
	webhookPath := "/webhook/" + randomHexString(RANDOM_SIZE)
	fmt.Printf("webhook path: %s\n", webhookPath)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.POST("/register", func(c echo.Context) error {
		reqTAInfo := new(struct {
			Domain        string
			GitRepository string
		})

		if c.Bind(reqTAInfo) != nil {
			return c.String(http.StatusBadRequest, "bad attestation")
		}

		taInfo := auditor.db.client.TAInfo.
			Create().
			SetDomain(reqTAInfo.Domain).
			SetGitRepository(reqTAInfo.GitRepository)

		_, err := taInfo.Save(*auditor.db.ctx)
		if err != nil {
			c.Error(err)
		}

		err = auditor.ct.Subscribe(reqTAInfo.Domain)
		if err != nil {
			c.Error(err)
		}

		return c.String(http.StatusOK, "ok")
	})

	e.POST("/compile", func(c echo.Context) error {
		idReq := new(struct {
			id int `json:"id"`
		})

		if c.Bind(idReq) != nil {
			return c.String(http.StatusBadRequest, "bad attestation")
		}

		taInfo, err := auditor.db.client.TAInfo.
			Query().Where(tainfo.IDEQ(idReq.id)).First(*auditor.db.ctx)

		if err != nil {
			c.Error(err)
		}

		commitId, uniqueId := compile(taInfo)

		taCode := auditor.db.client.TACode.
			Create().AddTaInfo(taInfo).SetCommitID(commitId).SetUniqueID(uniqueId)

		_, err = taCode.Save(*auditor.db.ctx)
		if err != nil {
			c.Error(err)
		}

		return c.String(http.StatusOK, "ok")
	})

	e.GET(webhookPath, func(c echo.Context) error {
		certs, err := auditor.ct.WebHookCertificates(c)
		if err != nil {
			c.Error(err)
		}

		err = auditor.AuditAll(certs)
		if err != nil {
			c.Error(err)
		}

		return c.String(http.StatusOK, "ok")
	})

	e.GET("/redirect", func(c echo.Context) error {
		back := c.Request().Header.Get("Referer")

		return c.Render(http.StatusOK, "redirect", back)
	})
}
