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
	"github.com/akakou/ra_webs/monitor"
	"github.com/akakou/ra_webs/monitor/ent"
	"github.com/akakou/ra_webs/monitor/ent/taserver"
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

var DefaultInterval = 10 * time.Minute

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

func (a *CrtshMonitor) Setup(monitor *monitor.Monitor) error {
	err := a.loadStream(monitor)
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

	a.loadFirstToClient()

	err = a.ctstream.Init()

	if err != nil {
		return err
	}

	return err
}

func (a *CrtshMonitor) PreCheck(domain string, exist bool, monitor *monitor.Monitor) error {
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

func (a *CrtshMonitor) Register(domain string, exist bool, monitor *monitor.Monitor) error {
	if exist {
		return nil
	}

	client, _, err := crtsh.AddByDomain(domain, a.ctstream.Client)
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

func (a *CrtshMonitor) updateInterval() {
	l := len(a.ctstream.Client.Clients)
	if l == 0 {
		l = 1
	}

	a.ctstream.Sleep = a.interval / time.Duration(l)
	a.ctstream.Client.Sleep = 0

	fmt.Printf("New Interval: %v \n", a.ctstream.Sleep)
}

func (a *CrtshMonitor) Run(monitor *monitor.Monitor) {
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

		serv, err := monitor.DB.Client.TAServer.
			Query().
			Order(ent.Desc(taserver.FieldID)).
			First(*monitor.DB.Ctx)

		fmt.Printf("[last] dbid: %v, domain: %v\n", serv.ID, domain)

		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		serv.Update().SetMonitorLogID(option.ID).Save(*monitor.DB.Ctx)

		err = Check(cert.PublicKey, serv)
		if err != nil {
			fmt.Printf("Violation: %v\n", err)
			revoke(serv, monitor)
			return
		}

		if !serv.IsActive {
			serv.Update().SetIsActive(true).Save(*monitor.DB.Ctx)
		}
	})
}
