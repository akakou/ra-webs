package monitor

import (
	"context"
	"errors"
	"time"

	crtapi "github.com/akakou/crtsh"
	"github.com/akakou/ctstream/monitor/crtsh"
	"github.com/akakou/ra_webs/verifier/core"
)

func NewCrtshMonitor(interval time.Duration, ctx context.Context) *CTMonitor[*crtsh.CrtshCTClient] {
	return &CTMonitor[*crtsh.CrtshCTClient]{
		ctx:      ctx,
		interval: interval,
		callback: CTMonitorCallbacks[*crtsh.CrtshCTClient]{
			defCTClient:  crtsh.NewCTClient,
			defCTsStream: crtsh.DefaultCTsStream,
			prepareFast:  loadFirstToCrthshClient,
			precheck:     preCheckWithCrtAPI,
		},
	}
}

func DefaultCrtshMonitor(ctx context.Context) *CTMonitor[*crtsh.CrtshCTClient] {
	return NewCrtshMonitor(DefaultInterval, ctx)
}

func loadFirstToCrthshClient(clients []*crtsh.CrtshCTClient, last int) {
	for _, c := range clients {
		c.ID = last
	}
}

func preCheckWithCrtAPI(domain string, verifier *core.Verifier) error {
	resp, err := crtapi.Fetch(domain, crtapi.EXCLUDE_EXPIRED)
	if err != nil {
		return err
	}

	if len(resp) != 0 {
		return errors.New(ERROR_FAILED_OTHER_CERTIFICATE_EXISTS)
	}

	return nil
}
