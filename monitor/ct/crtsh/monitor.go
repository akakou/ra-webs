package crtsh

import (
	"context"
	"fmt"
	"time"

	ctcore "github.com/akakou/ctstream/core"
	"github.com/akakou/ctstream/direct"
	"github.com/akakou/ctstream/monitor/crtsh"
	"github.com/akakou/ra-webs/monitor"
	ctx509 "github.com/google/certificate-transparency-go/x509"
)

type CrtshStream = ctcore.CTStream[*crtsh.CrtshCTClient]

type CrtshMonitor struct {
	CtStream *ctcore.CTStream[*crtsh.CrtshCTClient]
	Ctx      context.Context
	Last     int
	Interval time.Duration
	Monitor  *monitor.Monitor
}

var DefaultInterval = 10 * time.Minute

func New(domain string, interval time.Duration, monitor *monitor.Monitor, ctx context.Context) (*CrtshMonitor, error) {
	ctclient, err := crtsh.NewCTClient(domain)
	if err != nil {
		return nil, err
	}

	ctstream, err := ctcore.NewCTStream(ctclient, interval, ctx)
	if err != nil {
		return nil, err
	}

	return &CrtshMonitor{
		Ctx:      ctx,
		Interval: interval,
		Last:     0,
		CtStream: ctstream,
		Monitor:  monitor,
	}, nil
}

func Default(domain string, monitor *monitor.Monitor, ctx context.Context) (*CrtshMonitor, error) {
	return New(domain, DefaultInterval, monitor, ctx)
}

func (a *CrtshMonitor) Run() {
	a.CtStream.Run(func(cert *ctx509.Certificate, i int, params any, err error) {
		if err == nil {
		} else if err.Error() == direct.ERROR_FAILED_TO_NEW {
			return
		} else {
			fmt.Printf("Error: %v\n", err)
			return
		}

		option := params.(*crtsh.CrtshCTParams)
		fmt.Printf("[received] crtid: %v", option.ID)

		a.Monitor.Check(cert.PublicKey, option.ID)
	})
}
