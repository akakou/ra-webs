package monitor

import (
	"context"
	"errors"
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

type CTMonitorCallbacks[T ctcore.CtClient] struct {
	defCTClient  ctcore.DefaultCTClient[T]
	defCTsStream DefaultCTsStream[T]
	prepareFast  PrepareFastToCTClient[T]
}

type CTMonitor[T ctcore.CtClient] struct {
	ctstream   *ctcore.CTStream[*ctcore.CTClients[T]]
	ctx        context.Context
	lastLogger *goutils.File[int]
	last       int
	interval   time.Duration
	callback   CTMonitorCallbacks[T]
}

var DefaultInterval = 10 * time.Second

func NewCrtshMonitor(interval time.Duration, ctx context.Context) *CTMonitor[*crtsh.CrtshCTClient] {
	return &CTMonitor[*crtsh.CrtshCTClient]{
		ctx:      ctx,
		interval: interval,
		callback: CTMonitorCallbacks[*crtsh.CrtshCTClient]{
			defCTClient:  crtsh.NewCTClient,
			defCTsStream: crtsh.DefaultCTsStream,
			prepareFast:  loadFirstToCrthshClient,
		},
	}
}

func DefaultCrtshMonitor(ctx context.Context) *CTMonitor[*crtsh.CrtshCTClient] {
	return NewCrtshMonitor(DefaultInterval, ctx)
}

func NewSSLMateMonitor(interval time.Duration, ctx context.Context) *CTMonitor[*crtsh.CrtshCTClient] {
	return &CTMonitor[*crtsh.CrtshCTClient]{
		ctx:      ctx,
		interval: interval,
		callback: CTMonitorCallbacks[*crtsh.CrtshCTClient]{
			defCTClient:  crtsh.NewCTClient,
			defCTsStream: crtsh.DefaultCTsStream,
			prepareFast:  loadFirstToCrthshClient,
		},
	}
}

func DefaultCrtshMonitor(ctx context.Context) *CTMonitor[*crtsh.CrtshCTClient] {
	return NewCrtshMonitor(DefaultInterval, ctx)
}

func (a *CTMonitor[T]) Setup(verifier *core.Verifier) error {
	err := a.loadStream(verifier)
	if err != nil {
		return err
	}

	a.updateInterval()

	err = a.loadFileLogger()
	if err != nil {
		return err
	}

	err = a.loadFirst()
	if err != nil {
		return err
	}

	a.callback.prepareFast(a.ctstream.Client.Clients, a.last)

	err = a.ctstream.Init()

	if err != nil {
		return err
	}

	return err
}

func (a *CTMonitor[T]) PreCheck(domain string, exist bool, verifier *core.Verifier) error {
	if exist {
		return nil
	}

	resp, err := crtapi.Fetch(domain, crtapi.EXCLUDE_EXPIRED)
	if err != nil {
		return err
	}

	if len(resp) != 0 {
		return errors.New(ERROR_FAILED_OTHER_CERTIFICATE_EXISTS)
	}

	return nil

}

func (a *CTMonitor[T]) Register(domain string, exist bool, verifier *core.Verifier) error {
	if exist {
		return nil
	}

	client, _, err := ctcore.AddByDomain(domain, a.callback.defCTClient, a.ctstream.Client, nil)
	if err != nil {
		return err
	}

	a.updateInterval()

	err = client.Init()

	if err != nil {
		return err
	}

	return err
}

func (a *CTMonitor[T]) updateInterval() {
	l := len(a.ctstream.Client.Clients)
	if l == 0 {
		l = 1
	}

	a.ctstream.Sleep = a.interval / time.Duration(l)
	a.ctstream.Client.Sleep = 0

	fmt.Printf("New Interval: %v \n", a.ctstream.Sleep)
}

func (a *CTMonitor[T]) Run(verifier *core.Verifier) {
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
		if option.Client.Domain == lastClient.GetDomain() {
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

		serv.Update().SetMonitorLogID(option.ID).Save(*verifier.DB.Ctx)

		err = Check(cert.PublicKey, serv)
		if err != nil {
			fmt.Printf("Violation: %v\n", err)
			revoke(serv, verifier)
			return
		}

		if !serv.IsActive {
			serv.Update().SetIsActive(true).Save(*verifier.DB.Ctx)
		}
	})
}
