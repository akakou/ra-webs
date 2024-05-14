package ct

import (
	"crypto/x509"
	"fmt"

	"github.com/akakou/ra_webs/ttp/audit"
	"github.com/akakou/ra_webs/ttp/core"
	"github.com/akakou/sslmate-cert-search-api/api"
	"github.com/akakou/sslmate-cert-search-api/monitor"
	"github.com/labstack/echo/v4"
)

const LAST_FILE = "./last.txt"

type SSLMateCT struct {
	Monitors  monitor.Monitors
	Api       api.SSLMateSearchAPI
	BaseQuery api.Query
	Last      string
}

func NewSSLMateCT(token string) *SSLMateCT {
	ct := SSLMateCT{
		Monitors:  []monitor.Monitor{},
		Api:       *api.New(token),
		BaseQuery: api.Query{},
		Last:      "",
	}

	return &ct
}

func (ct *SSLMateCT) Update(ttp *core.TTP) error {
	last, err := readFile(LAST_FILE)
	if err != nil {
		return err
	}

	servers, err := ttp.DB.Client.TAServer.Query().All(*ttp.DB.Ctx)
	if err != nil {
		return err
	}

	monitors := []monitor.Monitor{}

	for _, server := range servers {
		query := ct.BaseQuery
		query.Domain = server.Domain
		query.After = last

		monitors = append(monitors, *monitor.New(&query, &ct.Api))
	}

	ct.Monitors = monitors

	return nil
}

func (ct *SSLMateCT) Setup(e *echo.Echo, ttp *core.TTP) error {
	go ct.Monitors.Run(func(certs []x509.Certificate, index *api.Index, err error) {
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

	})

	last, err := readFile(LAST_FILE)
	if err != nil {
		return err
	}

	ct.Last = last

	return nil
}

func (ct *SSLMateCT) Subscribe(string, ttp *core.TTP) error {
	return ct.Update(ttp)
}
