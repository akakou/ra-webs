package monitor

import (
	"bytes"
	"fmt"

	"github.com/cockroachdb/errors"

	"github.com/edgelesssys/ego/attestation"
	"github.com/edgelesssys/ego/attestation/tcbstatus"

	"github.com/akakou/ra-webs/devkit/core/attest"

	"github.com/akakou/ra-webs/monitor/builder"
	"github.com/akakou/ra-webs/monitor/serviceclient"
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

func CheckSourceHash(entry *serviceclient.EvidenceEntry, evidenceUniqueId []byte) error {
	uniqueId, err := builder.Build(entry.Repository, entry.CommitID)
	if err != nil {
		return errors.Wrap(errBuildFailed, err.Error())
	}

	fmt.Printf("unique id: '%x', '%x'\n", uniqueId, evidenceUniqueId)
	if !bytes.Equal(uniqueId, evidenceUniqueId) {
		return errUniqueIDInEvidenceMismatch
	}

	return nil
}
