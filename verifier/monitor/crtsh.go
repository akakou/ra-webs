package monitor

import (
	"context"
	"time"

	"github.com/akakou/ctstream/monitor/crtsh"
)

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

func loadFirstToCrthshClient(clients []*crtsh.CrtshCTClient, last int) {
	for _, c := range clients {
		c.ID = last
	}
}
