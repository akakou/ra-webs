package monitor

import (
	"fmt"

	"github.com/akakou/ctstream/monitor/crtsh"
	goutils "github.com/akakou/go-utils"
	"github.com/akakou/ra_webs/verifier/core"
	"github.com/akakou/ra_webs/verifier/ent/taserver"
)

const LOG_FILE_PATH = "first-log.txt"

const FILE_EMPLTY = "strconv.Atoi: parsing \"\": invalid syntax"

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

func (a *CrtshMonitor) loadFileLogger() error {
	lastLogger, err := goutils.OpenIntFile(LOG_FILE_PATH)
	if err != nil {
		return err
	}

	a.lastLogger = lastLogger

	return nil
}

func (a *CrtshMonitor) loadFirst() error {
	first, err := a.lastLogger.Restore()

	if err == nil {
	} else if err.Error() == FILE_EMPLTY {
		fmt.Printf("%v is empty !: %v\n", LOG_FILE_PATH, err)
		first = new(int)
		*first = 0
	} else {
		return err
	}

	fmt.Printf("First: %v\n", *first)
	return nil
}

func (a *CrtshMonitor) loadFirstToClient() {
	for _, c := range a.ctstream.Client.Clients {
		c.ID = a.last
	}
}
