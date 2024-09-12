package monitor

import (
	"github.com/akakou/ctstream/monitor/crtsh"
	"github.com/akakou/ra_webs/verifier/core"
	"github.com/akakou/ra_webs/verifier/ent/taserver"
)

func (a *CrtshMonitor) loadStream(verifier *core.Verifier) error {
	servers, err := verifier.DB.Client.TAServer.Query().Select(taserver.FieldDomain).All(*verifier.DB.Ctx)
	if err != nil {
		return err
	}

	domains := []string{}
	for _, serv := range servers {
		domains = append(domains, serv.Domain)
	}

	a.ctstream, err = crtsh.DefaultCTsStream(domains, a.ctx)
	if err != nil {
		return err
	}

	return nil
}
