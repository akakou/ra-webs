package sslmate

import (
	"context"

	"github.com/akakou/ctstream/core"
)

func NewCTsStream(clients *core.CTClients[*SSLMateCTClient], ctx context.Context) (*core.CTStream[*core.CTClients[*SSLMateCTClient]], error) {
	return core.NewCTStream(clients, 0, ctx)
}

func DefaultCTsStream(domains []string, ctx context.Context) (*core.CTStream[*core.CTClients[*SSLMateCTClient]], error) {
	client, err := DefaultCTClients(domains, ctx)
	if err != nil {
		return nil, err
	}

	stream, err := core.NewCTStream(client, core.DefaultEpochSleep, ctx)
	if err != nil {
		return nil, err
	}

	return stream, nil
}
