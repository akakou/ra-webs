package ttp

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"

	"github.com/akakou/ra_webs/core"
	"github.com/edgelesssys/ego/attestation"
)

func verifyAttestation(token string) (string, error) {
	signerId := flag.String("signer_id", "", "signer id")
	productId := flag.Uint64("product_id", 0, "product id")
	securityVersion := flag.Uint("security_version", 0, "security version")

	flag.Parse()

	signer, err := hex.DecodeString(*signerId)
	if err != nil {
		panic(err)
	}

	report, err := attestation.VerifyAzureAttestationToken(token, core.ATTEST_PROVIDER_URL)
	if err != nil {
		return "", err
	}

	if !bytes.Equal(report.SignerID, signer) {
		fmt.Printf("%v, %v", report.SignerID, signer)
		return "", errors.New("token does not contain the right signer id")
	}

	if binary.LittleEndian.Uint16(report.ProductID) != uint16(*productId) {
		return "", errors.New("token contains invalid product id")
	}

	if report.SecurityVersion < *securityVersion {
		return "", errors.New("token contains invalid security version number")
	}

	// Get certificate from the report.
	hashStrBytes := report.Data
	hashStr := string(hashStrBytes)

	return hashStr, nil
}
