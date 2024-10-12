package core

import (
	"time"
)

type CTClients[T CtClient] struct {
	Clients []T
	Sleep   time.Duration
}

func NewCTClients[T CtClient](clients []T, sleep time.Duration) (*CTClients[T], error) {
	return &CTClients[T]{
		Clients: clients,
		Sleep:   sleep,
	}, nil
}

func (clients *CTClients[T]) Init() error {
	for _, c := range clients.Clients {
		err := c.Init()
		if err != nil {
			return err
		}
	}

	return nil
}

func (clients *CTClients[T]) Next(callback Callback) {
	for _, c := range clients.Clients {
		c.Next(callback)
		time.Sleep(clients.Sleep)
	}
}

func (clients *CTClients[T]) GetDomain() string {
	return ""
}
