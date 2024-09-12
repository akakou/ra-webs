package monitor

import (
	"fmt"

	"github.com/akakou/ctstream/monitor/sslmate"
	goutils "github.com/akakou/go-utils"
	"github.com/akakou/ra_webs/verifier/core"
	"github.com/akakou/ra_webs/verifier/ent/taserver"
)

func (a *SSLMateMonitor) loadStream(verifier *core.Verifier) (*SSLMateStream, error) {
	servers, err := verifier.DB.Client.TAServer.Query().Select(taserver.FieldDomain).All(*verifier.DB.Ctx)
	if err != nil {
		return nil, err
	}

	domains := []string{}
	for _, serv := range servers {
		domains = append(domains, serv.Domain)
	}

	a.ctstream, err = sslmate.DefaultCTsStream(domains, a.ctx)
	if err != nil {
		return nil, err
	}

	return a.ctstream, nil
}

func (a *SSLMateMonitor) loadFileLogger() (*goutils.File[int], error) {
	lastLogger, err := goutils.OpenIntFile(LOG_FILE_PATH)
	if err != nil {
		return nil, err
	}

	a.lastLogger = lastLogger

	return lastLogger, nil
}

func (a *SSLMateMonitor) loadFirst(lastLogger *goutils.File[int]) (int, error) {
	first, err := lastLogger.Restore()

	if err == nil {
	} else if err.Error() == FILE_EMPLTY {
		fmt.Printf("%v is empty !: %v\n", LOG_FILE_PATH, err)
		first = new(int)
		*first = 0
	} else {
		return 0, err
	}

	fmt.Printf("First: %v\n", *first)
	return *first, nil
}
