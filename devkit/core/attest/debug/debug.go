package debug

import (
	"encoding/base64"
	"encoding/hex"

	"github.com/akakou/ra-webs/devkit/core/attest"
	"github.com/edgelesssys/ego/attestation"
	"github.com/edgelesssys/ego/attestation/tcbstatus"
)

var Debug = true

func debugAttestByAzure(data []byte) (string, error) {
	evidence := base64.StdEncoding.EncodeToString(data)
	return evidence, nil
}

func debugVerifyByAzure(evudence string) (*attestation.Report, error) {
	data, err := base64.StdEncoding.DecodeString(evudence)
	if err != nil {
		return nil, err
	}

	uniqueId, _ := hex.DecodeString("87d64f2dbe05ebabaf1993014078baa5a3c8d0089904a2013a070d479d6d13f6")
	return &attestation.Report{
		UniqueID:  uniqueId,
		Data:      data,
		Debug:     false,
		TCBStatus: tcbstatus.UpToDate,
	}, nil
}

const DEBUG_TOKEN = "this-is-ra-webs-debug-token-138484039348"

func EnableDebug() bool {
	Debug = true

	attest.Attest = debugAttestByAzure
	attest.Verify = debugVerifyByAzure

	return true
}
