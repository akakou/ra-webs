package ttp

import (
	"net/http"

	"strconv"

	"github.com/akakou/ra_webs/ttp/ent/ta"
	"github.com/labstack/echo/v4"
)

var registerTAApi = echoRoute{
	path: "/ta",
	f: func(auditor *Auditor) echoRouteFunc {
		return func(c echo.Context) error {
			reqTAInfo := new(struct {
				IP            string
				Domain        string
				GitRepository string
			})

			if c.Bind(reqTAInfo) != nil {
				return c.String(http.StatusBadRequest, "bad attestation")
			}

			taInfo := auditor.db.client.TA.
				Create().
				SetDomain(reqTAInfo.Domain).
				SetGit(reqTAInfo.GitRepository).
				SetIP(reqTAInfo.IP)

			_, err := taInfo.Save(*auditor.db.ctx)
			if err != nil {
				c.Error(err)
			}

			err = auditor.ct.Subscribe(reqTAInfo.Domain)
			if err != nil {
				c.Error(err)
			}

			return c.String(http.StatusOK, "ok")
		}
	},
}

var updateTAApi = echoRoute{
	path: "/ta/:id/update",
	f: func(auditor *Auditor) echoRouteFunc {
		return func(c echo.Context) error {
			idParam := c.Param("id")
			id, err := strconv.Atoi(idParam)
			if err != nil {
				return c.String(http.StatusBadRequest, "bad id")
			}

			taInfo, err := auditor.db.client.TA.
				Query().Where(ta.IDEQ(id)).First(*auditor.db.ctx)

			if err != nil {
				c.Error(err)
			}

			commitId, uniqueId := compile(taInfo)

			taCode := auditor.db.client.TACode.
				Create().AddTa(taInfo).SetCommitID(commitId).SetUniqueID(uniqueId)

			_, err = taCode.Save(*auditor.db.ctx)
			if err != nil {
				c.Error(err)
			}

			return c.String(http.StatusOK, "ok")
		}
	},
}
