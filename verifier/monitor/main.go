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
	last       int
	interval   time.Duration
}

var DefaultInterval = 10 * time.Second

func NewCrtshMonitor(interval time.Duration, ctx context.Context) *CrtshMonitor {
	return &CrtshMonitor{
		ctx:      ctx,
		interval: interval,
	}
}

func DefaultCrtshMonitor(ctx context.Context) *CrtshMonitor {
	return &CrtshMonitor{
		ctx:      ctx,
		interval: DefaultInterval,
	}
}

func (a *CrtshMonitor) Setup(verifier *core.Verifier) error {
	err := a.loadStream(verifier)
	if err != nil {
		return err
	}

	a.updateInteval()

	err = a.loadFileLogger()
	if err != nil {
		return err
	}

	err = a.loadFirst()
	if err != nil {
		return err
	}

	a.loadFirstToClient()

	err = a.ctstream.Init()

	if err != nil {
		return err
	}

	return err
}

func (a *CrtshMonitor) PreCheck(domain string, exist bool, verifier *core.Verifier) error {
	if exist {
		return nil
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

func (a *CrtshMonitor) Register(domain string, exist bool, verifier *core.Verifier) error {
	if exist {
		return nil
	}

	client, _, err := crtsh.AddByDomain(domain, a.ctstream.Client)
	if err != nil {
		return nil
	}

	a.updateInteval()

	err = client.Init()

	if err != nil {
		return nil
	}

	return err
}

func (a *CrtshMonitor) updateInteval() {
	l := len(a.ctstream.Client.Clients)
	if l == 0 {
		l = 1
	}

	a.ctstream.Sleep = a.interval / time.Duration(l)
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

		fmt.Printf("[received] crtid: %v, domain: %v\n", option.ID, domain)

		if option.ID > a.last {
			a.last = option.ID
		}

		clientsLen := len(a.ctstream.Client.Clients)
		lastClient := a.ctstream.Client.Clients[clientsLen-1]
		if option.Client.Domain == lastClient.Domain {
			err := a.lastLogger.Store(&a.last)
			if err != nil {
				panic(err)
			}
		}

		serv, err := verifier.DB.Client.TAServer.
			Query().
			Where(taserver.DomainEQ(domain)).
			Order(ent.Desc(taserver.FieldID)).
			First(*verifier.DB.Ctx)

		fmt.Printf("[last] dbid: %v, domain: %v\n", serv.ID, domain)

		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		err = Check(cert.PublicKey, serv)
		if err != nil {
			fmt.Printf("Violation: %v\n", err)
			revoke(option.ID, serv, verifier)
			return
		}

		if !serv.HasActivated {
			serv.Update().SetHasActivated(true).Save(*verifier.DB.Ctx)
		}
	})
}
