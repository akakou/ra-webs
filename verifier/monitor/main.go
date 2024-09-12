package monitor

import (
	"context"
	"fmt"
	"time"

	crtapi "github.com/akakou/crtsh"
	ctcore "github.com/akakou/ctstream/core"
	"github.com/akakou/ctstream/direct"
	"github.com/akakou/ctstream/monitor/crtsh"
	goutils "github.com/akakou/go-utils"
	"github.com/akakou/ra_webs/verifier/core"
	"github.com/akakou/ra_webs/verifier/ent"
	"github.com/akakou/ra_webs/verifier/ent/taserver"
	ctx509 "github.com/google/certificate-transparency-go/x509"
)

type CrtshStream = ctcore.CTStream[*ctcore.CTClients[*crtsh.CrtshCTClient]]

type CrtshMonitor struct {
	ctstream   *CrtshStream
	ctx        context.Context
	lastLogger *goutils.File[int]
	first      int
}

func NewCrtshMonitor(ctx context.Context) *CrtshMonitor {
	ctcore.DefaultEpochSleep = 10 * time.Second
	return &CrtshMonitor{
		ctx: ctx,
	}
}

func (a *CrtshMonitor) Setup(verifier *core.Verifier) error {
	err := a.loadStream(verifier)
	if err != nil {
		return err
	}

	err = a.loadFileLogger()
	if err != nil {
		return err
	}

	err = a.loadFirst()
	if err != nil {
		return err
	}

	return a.ctstream.Init()
}

func (a *CrtshMonitor) PreCheck(domain string, verifier *core.Verifier) error {
	_, _, err := crtsh.SelectByDomain(domain, a.ctstream.Client)

	if err == nil {
		return nil
	}

	if err.Error() != ctcore.ERROR_NOT_FOUND {
		return err
	}

	resp, err := crtapi.Fetch(domain, crtapi.EXCLUDE_EXPIRED)
	if err != nil {
		return err
	}

	if len(resp) != 0 {
		return fmt.Errorf(ERROR_FAILED_OTHER_CERTIFICATE_EXISTS)
	}

	return nil

}

func (a *CrtshMonitor) Register(domain string, verifier *core.Verifier) error {
	client, _, err := crtsh.AddByDomain(domain, a.ctstream.Client)
	if err != nil {
		return nil
	}

	err = client.Init()

	if err != nil {
		return nil
	}

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

		if option.ID > a.first {
			a.first = option.ID
			a.lastLogger.Store(&a.first)
		}

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
