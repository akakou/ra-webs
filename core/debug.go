package core

import "github.com/edgelesssys/ego/attestation"

func debugAttestByAzure(data []byte) (string, error) {
	return "", nil
}

func debugVerifyByAzure(token string, data []byte) (*attestation.Report, error) {
	return &attestation.Report{
		UniqueID: []byte{1, 2, 3},
		Data:     []byte{4, 5, 6},
	}, nil
}

func EnableDebug() {
	DEBUG = true

	AttestByAzure = debugAttestByAzure
	VerifyByAzure = debugVerifyByAzure
}
