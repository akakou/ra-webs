package monitor

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"fmt"

	"github.com/cockroachdb/errors"

	"github.com/edgelesssys/ego/attestation"
	"github.com/edgelesssys/ego/attestation/tcbstatus"

	core "github.com/akakou/ra-webs/core"
	"github.com/akakou/ra-webs/log/api/interfacestruct"
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

func CheckEvidence(quote string) (*attestation.Report, error) {
	report, err := core.Verify(quote)
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

func CheckPublicKey(ctPublicKey publicKey, logPublicKey []byte) error {
	unmarshaledPublicKey, isRSA := ctPublicKey.(*rsa.PublicKey)
	if !isRSA {
		return errPublicKeyIsNotRSA
	}

	ctPublicKeyBuf := x509.MarshalPKCS1PublicKey(unmarshaledPublicKey)
	fmt.Printf("compireing public key:\n%v\n!=%v\n\n", ctPublicKeyBuf, logPublicKey)

	if !bytes.Equal(ctPublicKeyBuf, logPublicKey) {
		return errPublicKeyNotMatched
	}

	return nil
}

func CheckSourceHash(log *interfacestruct.TA, evidenceUniqueId []byte) error {
	uniqueId, err := builder.Build(log.Repository, log.CommitID)
	if err != nil {
		return errors.Wrap(errBuildFailed, err.Error())
	}

	if !bytes.Equal(uniqueId, evidenceUniqueId) {
		return errUniqueIDInEvidenceMismatch
	}

	return nil
}
