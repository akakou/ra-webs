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
	MaxSleep  time.Duration
	// Last      string
}

func Nop(c []x509.Certificate, i *api.Index, err error) {}

func DefaultSSLMateCT(token string) *SSLMateCT {
	ct := SSLMateCT{
		Monitors: monitor.Monitors{
			Monitors: []monitor.Monitor{},
		},
		Api:       *api.New(token),
		BaseQuery: monitor.DefaultQuery,
		// Last:      "",
		MaxSleep: DEFAULT_MAX_SLEEP,
	}

	return &ct
}

func (ct *SSLMateCT) Setup(e *echo.Echo, ttp *core.TTP) error {
	// last, err := readFile(LAST_FILE)
	// if err != nil {
	// 	return err
	// }

	err := ct.LoadFromDB(ttp)
	if err != nil {
		return err
	}

	ct.Monitors.Next(Nop)

	go ct.Monitors.Loop(Routine(ttp))

	fmt.Println("ct started...")

	return nil
}

func Routine(ttp *core.TTP) monitor.Callback {
	return func(certs []x509.Certificate, index *api.Index, err error) {
		fmt.Println("Now CT check running...")

		if err != nil {
			fmt.Printf("ct error: %v\n", err)
		}

		err = audit.AuditAll(ttp, certs)
		if err != nil {
			fmt.Printf("ct error: %v\n", err)
		}

		// err = writeFile(LAST_FILE, index.Last)
		// if err != nil {
		// 	fmt.Printf("ct error: %v\n", err)
		// }
	}
}

func (ct *SSLMateCT) LoadFromDB(ttp *core.TTP) error {
	fmt.Println("starting sync from db...")
	monitors := []monitor.Monitor{}

	domains, err := ttp.DB.Client.TAServer.Query().Select(taserver.FieldDomain).Strings(*ttp.DB.Ctx)

	if err != nil {
		return err
	}

	domains = removeDeplication(domains)

	for _, domain := range domains {
		query := ct.BaseQuery
		query.Domain = domain
		query.After = ""

		m := monitor.Monitor{
			Query: &query,
			Api:   &ct.Api,
			Sleep: DEFAULT_MAX_SLEEP,
		}

		monitors = append(monitors, m)
	}

	ct.Monitors.Monitors = monitors
	ct.ResetSleep()

	return nil
}

func (ct *SSLMateCT) ResetSleep() time.Duration {
	sleep := ct.MaxSleep

	if len(ct.Monitors.Monitors) > 0 {
		sleep /= time.Duration(len(ct.Monitors.Monitors))
	}

	for _, monitor := range ct.Monitors.Monitors {
		monitor.Sleep = sleep
	}

	return sleep
}

func hasMonitorForDomain(domain string, monitors monitor.Monitors) bool {
	hasMonitor := false
	for _, monitor := range monitors.Monitors {
		hasMonitor = hasMonitor || (monitor.Query.Domain == domain)
	}

	return hasMonitor
}

func (ct *SSLMateCT) Insert(domain string) error {
	if hasMonitorForDomain(domain, ct.Monitors) {
		return nil
	}

	sleep := ct.ResetSleep()

	m := monitor.DefaultMonitor(domain)
	m.Sleep = sleep

	_, _, err := m.Next()
	if err != nil {
		return err
	}

	ct.Monitors.Monitors = append(ct.Monitors.Monitors, *m)

	return nil
}

func (ct *SSLMateCT) Subscribe(domain string, ttp *core.TTP) error {
	return ct.Insert(domain)
}
