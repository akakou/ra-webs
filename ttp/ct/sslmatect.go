package ct

import (
	"crypto/x509"
	"fmt"
	"time"

	"github.com/akakou/ra_webs/ttp/audit"
	"github.com/akakou/ra_webs/ttp/core"
	"github.com/akakou/ra_webs/ttp/ent/taserver"
	"github.com/akakou/sslmate-cert-search-api/api"
	"github.com/akakou/sslmate-cert-search-api/monitor"
	"github.com/labstack/echo/v4"
)

const LAST_FILE = "./last.txt"
const DEFAULT_MAX_SLEEP = monitor.DEFAULT_SLEEP

type SSLMateCT struct {
	Monitors  monitor.Monitors
	Api       api.SSLMateSearchAPI
	BaseQuery api.Query
	Last      string
	Sleep     time.Duration
}

func NewSSLMateCT(token string) *SSLMateCT {
	ct := SSLMateCT{
		Monitors:  []monitor.Monitor{},
		Api:       *api.New(token),
		BaseQuery: api.Query{},
		Last:      "",
		Sleep:     DEFAULT_MAX_SLEEP,
	}

	return &ct
}

func (ct *SSLMateCT) Setup(e *echo.Echo, ttp *core.TTP) error {
	last, err := readFile(LAST_FILE)
	if err != nil {
		return err
	}

	ct.Last = last

	err = ct.SyncFromDB(ttp)
	if err != nil {
		return err
	}

	go ct.Monitors.Run(func(certs []x509.Certificate, index *api.Index, err error) {
		fmt.Println("Now CT check running...")

		if err != nil {
			fmt.Printf("%v", err)
		}

		err = audit.AuditAll(ttp, certs)
		if err != nil {
			fmt.Printf("%v", err)
		}

		err = writeFile(index.Last)
		if err != nil {
			fmt.Printf("%v", err)
		}

		time.Sleep(ct.Sleep)
	})

	fmt.Println("ct started...")

	return nil
}

func (ct *SSLMateCT) SyncFromDB(ttp *core.TTP) error {
	monitors := []monitor.Monitor{}

	domains, err := ttp.DB.Client.TAServer.Query().Select(taserver.FieldDomain).Strings(*ttp.DB.Ctx)

	if err != nil {
		return err
	}

	domains = removeDeplication(domains)

	for _, domain := range domains {
		query := ct.BaseQuery
		query.Domain = domain
		query.After = ct.Last

		m := monitor.Monitor{
			Query: &query,
			Api:   &ct.Api,
			Sleep: 0,
		}

		monitors = append(monitors, m)
	}

	ct.Monitors = monitors
	ct.Sleep = DEFAULT_MAX_SLEEP / time.Duration(len(domains))

	return nil
}

func (ct *SSLMateCT) Subscribe(_ string, ttp *core.TTP) error {
	return ct.SyncFromDB(ttp)
}
