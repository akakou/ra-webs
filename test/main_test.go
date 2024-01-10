package test

import (
	"testing"

	"github.com/akakou/ra_webs/ta"
	"github.com/akakou/ra_webs/ttp"
	"github.com/stretchr/testify/assert"
)

func TestProvisioning(t *testing.T) {
	ttpPort := ":12347"
	ta.SCHEME = "http://"

	e := ttp.DefaultTTPServer()
	go e.Start(ttpPort)

	config := ta.RAConfig{
		Domain:    "test",
		TTPDomain: "localhost" + ttpPort,
	}

	_, err := config.Provisioning()
	assert.NoError(t, err)

	e.Server.Close()
}
