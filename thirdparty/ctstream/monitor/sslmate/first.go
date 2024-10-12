package sslmate

import (
	"strconv"

	"github.com/akakou/ctstream/core"
)

func GetFirst(ctClients *core.CTClients[*SSLMateCTClient]) int {
	var first int

	for _, client := range ctClients.Clients {
		tmp, _ := strconv.Atoi(client.First)

		if tmp > first {
			first = tmp
		}
	}

	return first
}

func SetFirst(first int, ctClients *core.CTClients[*SSLMateCTClient]) {
	for _, client := range ctClients.Clients {
		client.First = strconv.Itoa(first)
	}
}
