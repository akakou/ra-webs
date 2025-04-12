package debug

import (
	"github.com/akakou/ra-webs/core/attest"
	"github.com/edgelesssys/ego/attestation"
	"github.com/edgelesssys/ego/attestation/tcbstatus"
)

var Debug = false

func debugAttestByAzure(data []byte) (string, error) {
	return "this is evidence", nil
}

func debugVerifyByAzure(evudence string) (*attestation.Report, error) {
	return &attestation.Report{
		UniqueID:  []byte{1, 2, 3},
		Data:      []byte{4, 5, 6},
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
