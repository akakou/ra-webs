package monitor

import (
	"bytes"
	"crypto/rsa"
	"fmt"

	"github.com/cockroachdb/errors"

	"github.com/edgelesssys/ego/attestation"
	"github.com/edgelesssys/ego/attestation/tcbstatus"

	"github.com/akakou/ra-webs/core/attest"
	"github.com/akakou/ra-webs/core/sign"

	"github.com/akakou/ra-webs/log/api/io"
	"github.com/akakou/ra-webs/monitor/builder"
)

var (
	errEnclaveIsDebugMode          = errors.New("enclave is in debug mode")
	errTCBStatusNotUpToDate        = errors.New("TCB status is not up to date")
	errUniqueIDInEvidenceMismatch  = errors.New("unique ID in evidence mismatch")
	errUniqueIDMadeByBuildMismatch = errors.New("unique ID made by build mismatch")
	errPublicKeyNotMatched         = errors.New("public key not matched")
	errPublicKeyIsNotRSA           = errors.New("public key is not RSA")
	errBuildFailed                 = errors.New("build failed")
)

func CheckSignature(signature []byte, l *sign.LogPlain, publicKey *rsa.PublicKey) error {
	err := sign.Verify(signature, l, publicKey)
	return err

}

func CheckEvidence(evidence string) (*attestation.Report, error) {
	report, err := attest.Verify(evidence)
	if err != nil {
		return &attestation.Report{
			UniqueID:  []byte{},
			Debug:     false,
			TCBStatus: tcbstatus.Unknown,
		}, err
	}

	if report.Debug {
		return &attestation.Report{
			UniqueID:  []byte{},
			Debug:     true,
			TCBStatus: tcbstatus.Unknown,
		}, errEnclaveIsDebugMode
	}

	if report.TCBStatus != tcbstatus.UpToDate {
		return &attestation.Report{
			UniqueID:  []byte{},
			Debug:     false,
			TCBStatus: tcbstatus.UpToDate,
		}, errTCBStatusNotUpToDate
	}

	return report, nil
}

func CheckSourceHash(log *io.TA, evidenceUniqueId []byte) error {
	uniqueId, err := builder.Build(log.Repository, log.CommitID)
	if err != nil {
		return errors.Wrap(errBuildFailed, err.Error())
	}

	fmt.Printf("unique id: '%x', '%x'\n", uniqueId, evidenceUniqueId)
	if !bytes.Equal(uniqueId, evidenceUniqueId) {
		return errUniqueIDInEvidenceMismatch
	}

	return nil
}
