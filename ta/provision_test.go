package ta

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRAConfiggenerateKeyPair(t *testing.T) {
	raConfig := RAConfig{
		TTPDomain: "ttp.example.com",
		Domain:    "ta.example.com",
	}

	ra := NewRA(&raConfig)

	_, _, err := ra.generateKeyPair()

	assert.NoError(t, err)
}
