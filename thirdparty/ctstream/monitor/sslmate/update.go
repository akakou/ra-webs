package sslmate

import (
	"github.com/akakou/ctstream/core"
)

func null() *SSLMateCTClient { return nil }

func SelectByDomain(
	domain string,
	clients *core.CTClients[*SSLMateCTClient],
) (*SSLMateCTClient, int, error) {
	client, i, err := core.SelectByDomain(domain, clients, null)
	return client, i, err
}

func AddByDomain(
	domain string,
	clients *core.CTClients[*SSLMateCTClient],
) (*SSLMateCTClient, int, error) {
	client, i, err := core.AddByDomain(domain, DefaultCTClient, clients, null)
	return client, i, err
}

func DelByDomain(
	domain string,
	clients *core.CTClients[*SSLMateCTClient],
) (*SSLMateCTClient, int, error) {
	client, i, err := core.DelByDomain(domain, clients, null)
	return client, i, err
}
