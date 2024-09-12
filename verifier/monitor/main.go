package monitor

import (
	"context"
	"fmt"

	ctcore "github.com/akakou/ctstream/core"
	"github.com/akakou/ctstream/direct"
	"github.com/akakou/ctstream/monitor/sslmate"
	goutils "github.com/akakou/go-utils"
	"github.com/akakou/ra_webs/verifier/core"
	"github.com/akakou/ra_webs/verifier/ent"
	"github.com/akakou/ra_webs/verifier/ent/taserver"
	ctx509 "github.com/google/certificate-transparency-go/x509"
)

const LOG_FILE_PATH = "last.log"
const FILE_EMPLTY = "strconv.Atoi: parsing \"\": invalid syntax"

type SSLMateStream = ctcore.ConcurrentCTsStream[*ctcore.CTStream[*sslmate.SSLMateCTClient]]

type SSLMateMonitor struct {
	ctstream   *SSLMateStream
	lastLogger *goutils.File[int]
	ctx        context.Context
}

func NewSSLMateMonitor(ctx context.Context) *SSLMateMonitor {
	return &SSLMateMonitor{
		ctx: ctx,
	}

}

func (a *SSLMateMonitor) Setup(verifier *core.Verifier) error {
	stream, err := a.loadStream(verifier)
	if err != nil {
		return err
	}

	a.ctstream = stream

	lastLogger, err := a.loadFileLogger(verifier)
	if err != nil {
		return err
	}

	a.lastLogger = lastLogger

	first, err := a.loadFirst(lastLogger)
	if err != nil {
		return err
	}

	sslmate.SetFirst(first, a.ctstream)

	return a.ctstream.Init()
}

func (a *SSLMateMonitor) Register(domain string, verifier *core.Verifier) error {
	err := sslmate.AddByDomain(domain, context.Background(), a.ctstream)
	return err
}

func (a *SSLMateMonitor) Run(verifier *core.Verifier) {
	a.ctstream.Run(func(cert *ctx509.Certificate, i int, params any, err error) {
		if err == nil {
		} else if err.Error() == direct.ERROR_FAILED_TO_NEW {
			return
		} else {
			fmt.Printf("Error: %v\n", err)
			return
		}

		option := params.(sslmate.SSLMateCTParams)
		domain := option.Client.Domain

		serv, err := verifier.DB.Client.TAServer.
			Query().
			Where(taserver.DomainEQ(domain)).
			Order(ent.Desc(taserver.FieldID)).
			First(*verifier.DB.Ctx)

		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		last := sslmate.GetFirst(a.ctstream)

		err = a.lastLogger.Store(&last)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}

		err = Check(cert.PublicKey, serv)
		if err != nil {
			fmt.Printf("Violation: %v\n", err)
			revoke(serv, verifier)
			return
		}

		if !serv.HasActivated {
			serv.Update().SetHasActivated(true).Save(*verifier.DB.Ctx)
		}
	})
}
