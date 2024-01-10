package ta

import (
	"crypto/tls"
	"log"

	"github.com/edgelesssys/ego/enclave"
)

const ATTEST_PROVIDER_URL = "https://shareduks.uks.attest.azure.net"

func AttestateByAzure(tlsConfig *tls.Config) string {
	token, err := enclave.CreateAzureAttestationToken(tlsConfig.Certificates[0].Certificate[0], ATTEST_PROVIDER_URL)
	if err != nil {
		log.Print("Run without attestation!!!!!\n")
	} else {
		log.Print("Created an Microsoft Azure Attestation Token.")
	}

	return token
}
