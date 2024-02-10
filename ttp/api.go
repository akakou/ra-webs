package ttp

import (
	"net/http"

	"github.com/akakou/ra_webs/ttp/ent/tainfo"
	"github.com/labstack/echo/v4"
)

func (auditServ *AuditServer) register() echoRoute {
	return func(c echo.Context) error {
		reqTAInfo := new(struct {
			Domain        string
			GitRepository string
		})

		if c.Bind(reqTAInfo) != nil {
			return c.String(http.StatusBadRequest, "bad attestation")
		}

		taInfo := auditServ.auditor.db.client.TAInfo.
			Create().
			SetDomain(reqTAInfo.Domain).
			SetGitRepository(reqTAInfo.GitRepository)

		_, err := taInfo.Save(*auditServ.auditor.db.ctx)
		if err != nil {
			c.Error(err)
		}

		err = auditServ.auditor.ct.Subscribe(reqTAInfo.Domain)
		if err != nil {
			c.Error(err)
		}

		return c.String(http.StatusOK, "ok")
	}
}

func (auditServ *AuditServer) compile() echoRoute {
	return func(c echo.Context) error {
		idReq := new(struct {
			Id int `json:"id"`
		})

		if c.Bind(idReq) != nil {
			return c.String(http.StatusBadRequest, "bad attestation")
		}

		taInfo, err := auditServ.auditor.db.client.TAInfo.
			Query().Where(tainfo.IDEQ(idReq.Id)).First(*auditServ.auditor.db.ctx)

		if err != nil {
			c.Error(err)
		}

		commitId, uniqueId := compile(taInfo)

		taCode := auditServ.auditor.db.client.TACode.
			Create().AddTaInfo(taInfo).SetCommitID(commitId).SetUniqueID(uniqueId)

		_, err = taCode.Save(*auditServ.auditor.db.ctx)
		if err != nil {
			c.Error(err)
		}

		return c.String(http.StatusOK, "ok")
	}
}
