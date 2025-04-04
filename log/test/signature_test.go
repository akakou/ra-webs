package test

import (
	"testing"

	"github.com/akakou/ra-webs/log/core"
	"github.com/stretchr/testify/assert"
)

func testSignature(t *testing.T, log *core.Log) {
	signature, err := core.Sign(log, &reqData[0])
	assert.NoError(t, err, "Sigining failed")

	_, err = core.Verify(signature, log.Domain, &reqData[0], log.VerifyKey)
	assert.NoError(t, err, "Signature verification failed")
}
