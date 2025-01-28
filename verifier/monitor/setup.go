package monitor

import (
	"context"
	"fmt"

	ctcore "github.com/akakou/ctstream/core"
	goutils "github.com/akakou/go-utils"

	"github.com/akakou/ra_webs/verifier/core"
	"github.com/akakou/ra_webs/verifier/ent/taserver"
)

const LOG_FILE_PATH = "last-log.txt"

const FILE_EMPLTY = "strconv.Atoi: parsing \"\": invalid syntax"

type DefaultCTsStream[T ctcore.CtClient] func([]string, context.Context) (*ctcore.CTStream[*ctcore.CTClients[T]], error)
type PrepareFastToCTClient[T ctcore.CtClient] func([]T, int)

func (a *CTMonitor[T]) loadStream(verifier *core.Verifier) error {
	servers, err := verifier.DB.Client.TAServer.Query().Select(taserver.FieldDomain).All(*verifier.DB.Ctx)
	if err != nil {
		return fmt.Errorf("%v:%v", ERROR_SELECT_TAS, err)
	}

	domains := []string{}
	for _, serv := range servers {
		domains = append(domains, serv.Domain)
	}

	a.ctstream, err = a.callback.defCTsStream(domains, a.ctx)
	if err != nil {
		return fmt.Errorf("%v:%v", ERROR_FAILED_TO_NEW_CTSSTREAM, err)
	}

	return nil
}

func (a *CTMonitor[T]) loadFileLogger() error {
	lastLogger, err := goutils.OpenIntFile(LOG_FILE_PATH)
	if err != nil {
		return err
	}

	a.lastLogger = lastLogger

	return nil
}

func (a *CTMonitor[T]) loadFirst() error {
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
