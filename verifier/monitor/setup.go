package monitor

import (
	"fmt"

	"github.com/akakou/ctstream/monitor/crtsh"
	"github.com/akakou/ra_webs/verifier/core"
	"github.com/akakou/ra_webs/verifier/ent/taserver"
)

func (a *CrtshMonitor) loadStream(verifier *core.Verifier) error {
	servers, err := verifier.DB.Client.TAServer.Query().Select(taserver.FieldDomain).All(*verifier.DB.Ctx)
	if err != nil {
		return fmt.Errorf("%v:%v", ERROR_SELECT_TAS, err)
	}

	domains := []string{}
	for _, serv := range servers {
		domains = append(domains, serv.Domain)
	}

	a.ctstream, err = crtsh.DefaultCTsStream(domains, a.ctx)
	if err != nil {
		return fmt.Errorf("%v:%v", ERROR_FAILED_TO_NEW_CTSSTREAM, err)
	}

	return nil
}
