package monitor

import (
	"context"
	"fmt"

	ctcore "github.com/akakou/ctstream/core"
	"github.com/akakou/ctstream/direct"
	"github.com/akakou/ctstream/monitor/crtsh"
	"github.com/akakou/ra_webs/verifier/core"
	"github.com/akakou/ra_webs/verifier/ent"
	"github.com/akakou/ra_webs/verifier/ent/taserver"
	ctx509 "github.com/google/certificate-transparency-go/x509"
)

type CrtshStream = ctcore.ConcurrentCTsStream[*ctcore.CTStream[*crtsh.CrtshCTClient]]

type CrtshMonitor struct {
	ctstream *CrtshStream
	ctx      context.Context
}

func NewCrtshMonitor(ctx context.Context) *CrtshMonitor {
	return &CrtshMonitor{
		ctx: ctx,
	}
}

func (a *CrtshMonitor) Setup(verifier *core.Verifier) error {
	a.loadStream(verifier)
	return a.ctstream.Init()
}

func (a *CrtshMonitor) Register(domain string, verifier *core.Verifier) error {
	err := crtsh.AddByDomain(domain, context.Background(), a.ctstream)
	return err
}

func (a *CrtshMonitor) Run(verifier *core.Verifier) {
	a.ctstream.Run(func(cert *ctx509.Certificate, i int, params any, err error) {
		if err == nil {
		} else if err.Error() == direct.ERROR_FAILED_TO_NEW {
			return
		} else {
			fmt.Printf("Error: %v\n", err)
			return
		}

		option := params.(*crtsh.CrtshCTParams)
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
