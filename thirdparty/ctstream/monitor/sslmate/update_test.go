package sslmate

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdate(t *testing.T) {
	domain := "test.example.com"

	clients, _ := DefaultCTClients([]string{}, context.Background())

	t.Run("TestAdd", func(t *testing.T) {
		const domain = "test.example.com"
		_, i, err := AddByDomain(domain, clients)
		assert.NoError(t, err)
		assert.Equal(t, i, 0)
		assert.Equal(t, len(clients.Clients), 1)
	})

	t.Run("TestDel", func(t *testing.T) {
		_, i, err := DelByDomain(domain, clients)
		assert.NoError(t, err)
		assert.Equal(t, i, 0)
		assert.Equal(t, len(clients.Clients), 0)
	})

}
