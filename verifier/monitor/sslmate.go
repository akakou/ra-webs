package monitor

import (
	"context"
	"errors"
	"time"

	"github.com/akakou/ctstream/monitor/sslmate"
	"github.com/akakou/ra_webs/verifier/core"
	sslmateapi "github.com/akakou/sslmate-cert-search-api/api"
)

func NewSSLMateMonitor(interval time.Duration, ctx context.Context) *CTMonitor[*sslmate.SSLMateCTClient] {
	return &CTMonitor[*sslmate.SSLMateCTClient]{
		ctx:      ctx,
		interval: interval,
		callback: CTMonitorCallbacks[*sslmate.SSLMateCTClient]{
			defCTClient:  sslmate.DefaultCTClient,
			defCTsStream: sslmate.DefaultCTsStream,
			prepareFast:  loadFirstToSSLMateClient,
			precheck:     preCheckWithSSLMate,
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

func preCheckWithSSLMate(domain string, verifier *core.Verifier) error {
	api := sslmateapi.Default()

	certs, _, err := api.Search(&sslmateapi.Query{
		Domain: domain,
	})

	if len(certs) != 0 {
		return errors.New(ERROR_FAILED_OTHER_CERTIFICATE_EXISTS)
	}

	return err
}
