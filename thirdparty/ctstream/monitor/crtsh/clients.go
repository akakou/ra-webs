package crtsh

import (
	"context"
	"time"

	"github.com/akakou/ctstream/core"
)

func NewCTClients(clients []*CrtshCTClient, sleep time.Duration) (*core.CTClients[*CrtshCTClient], error) {
	res, err := core.NewCTClients(clients, sleep)
	return res, err
}

func DefaultCTClients(domains []string, context context.Context) (*core.CTClients[*CrtshCTClient], error) {
	var clients []*CrtshCTClient

	for _, domain := range domains {
		client, err := NewCTClient(domain)
		if err != nil {
			return nil, err
		}

		clients = append(clients, client)
	}

	res, err := core.NewCTClients(clients, core.DefaultEpochSleep)

	return res, err
}
