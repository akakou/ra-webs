package sslmate

import (
	"context"
	"time"

	"github.com/akakou/ctstream/core"
)

func NewCTClients(clients []*SSLMateCTClient, sleep time.Duration) (*core.CTClients[*SSLMateCTClient], error) {
	res, err := core.NewCTClients(clients, sleep)
	return res, err
}

func DefaultCTClients(domains []string, context context.Context) (*core.CTClients[*SSLMateCTClient], error) {
	var clients []*SSLMateCTClient

	for _, domain := range domains {
		client, err := DefaultCTClient(domain)
		if err != nil {
			return nil, err
		}

		clients = append(clients, client)
	}

	res, err := core.NewCTClients(clients, core.DefaultEpochSleep)

	return res, err
}
