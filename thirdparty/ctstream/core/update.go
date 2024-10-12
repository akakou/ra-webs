package core

import (
	"errors"
)

const ERROR_NOT_FOUND = "not found"

type DefaultCTClient[C CtClient] func(string) (C, error)

func SelectByDomain[C CtClient](
	domain string,
	clients *CTClients[C],
	null func() C,
) (C, int, error) {
	for i, c := range clients.Clients {
		if c.GetDomain() == domain {
			return c, i, nil
		}
	}

	return null(), -1, errors.New(ERROR_NOT_FOUND)
}

func AddByDomain[C CtClient](
	domain string,
	def DefaultCTClient[C],
	clients *CTClients[C],
	null func() C,
) (C, int, error) {
	client, err := def(domain)

	if err != nil {
		return null(), 0, err
	}

	clients.Clients = append(clients.Clients, client)

	return client, len(clients.Clients) - 1, nil
}

func DelByDomain[C CtClient](
	domain string,
	clients *CTClients[C],
	null func() C,
) (C, int, error) {
	client, i, err := SelectByDomain(domain, clients, null)
	if err != nil {
		return null(), 0, err
	}

	clients.Clients = append(clients.Clients[:i], clients.Clients[i+1:]...)

	return client, i, nil
}
