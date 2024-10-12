package crtsh

import (
	"github.com/akakou/crtsh"
	"github.com/akakou/ctstream/core"
	"github.com/akakou/ctstream/utils"
)

type CrtshCTClient struct {
	ID     int
	Domain string
}

type CrtshCTParams struct {
	ID     int
	Client *CrtshCTClient
}

func NewCTClient(domain string) (*CrtshCTClient, error) {
	return &CrtshCTClient{
		ID:     0,
		Domain: domain,
	}, nil
}

func (client *CrtshCTClient) Init() error {
	_, err := crtsh.Fetch(client.Domain, crtsh.EXCLUDE_EXPIRED)
	return err
}

func (client *CrtshCTClient) Next(callback core.Callback) {
	entries, err := crtsh.Fetch(client.Domain, crtsh.EXCLUDE_EXPIRED)
	if err != nil {
		callback(nil, 0, &CrtshCTParams{}, err)
		return
	}

	for i, entry := range entries {
		if entry.ID <= client.ID {
			continue
		}

		c, err := utils.ReformatCertificate(entry.Certificate)
		callback(c, i, &CrtshCTParams{
			ID:     entries[i].ID,
			Client: client,
		}, err)

		client.ID = entry.ID
	}
}

func (client *CrtshCTClient) GetDomain() string {
	return client.Domain
}
