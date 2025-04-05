package core

import "github.com/edgelesssys/ego/attestation"

func debugAttestByAzure(data []byte) (string, error) {
	return "", nil
}

func debugVerifyByAzure(evudence string) (*attestation.Report, error) {
	return &attestation.Report{
		UniqueID: []byte{1, 2, 3},
		Data:     []byte{4, 5, 6},
	}, nil
}

const DEBUG_TOKEN = "this-is-ra-webs-debug-token-138484039348"

func EnableDebug() {
	DEBUG = true

	Attest = debugAttestByAzure
	Verify = debugVerifyByAzure
}
