package monitor

import (
	"context"
	"time"

	"github.com/akakou/ctstream/monitor/sslmate"
)

func NewSSLMateMonitor(interval time.Duration, ctx context.Context) *CTMonitor[*sslmate.SSLMateCTClient] {
	return &CTMonitor[*sslmate.SSLMateCTClient]{
		ctx:      ctx,
		interval: interval,
		callback: CTMonitorCallbacks[*sslmate.SSLMateCTClient]{
			defCTClient:  sslmate.DefaultCTClient,
			defCTsStream: sslmate.DefaultCTsStream,
			prepareFast:  loadFirstToSSLMateClient,
		},
	}
}

func DefaultSSLMateMonitor(ctx context.Context) *CTMonitor[*sslmate.SSLMateCTClient] {
	return NewSSLMateMonitor(DefaultInterval, ctx)
}

func loadFirstToSSLMateClient(clients []*sslmate.SSLMateCTClient, last int) {
	for _, c := range clients {
		c.First = string(last)
	}
}
