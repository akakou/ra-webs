package monitor

import (
	"crypto/rsa"
	"crypto/x509"
	"fmt"

	"github.com/akakou/crtsh"
	"github.com/akakou/ra-webs/core/sign"
	"github.com/akakou/ra-webs/log/api/io"
	"github.com/akakou/ra-webs/monitor/ent/atlog"
	"github.com/akakou/ra-webs/monitor/ent/ctlog"
)

type publicKey interface{}

func (monitor *Monitor) Monitor(ctLogs []crtsh.CertificateEntry) {
	atLogs, err := monitor.LogClient.Fetch()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	for _, entry := range atLogs {
		monitor.MonitorATLog(entry)
	}

	for _, entry := range ctLogs {
		monitor.MonitorCTLog(entry)
	}
}

func (monitor *Monitor) MonitorATLog(log *io.TA) {
	var err error

	exist := monitor.DB.Client.ATLog.Query().
		Where(atlog.EvidenceEQ(log.Evidence)).
		ExistX(*monitor.DB.Ctx)

	if exist {
		return
	}

	report, err := CheckEvidence(log.Evidence)
	if err != nil {
		fmt.Printf("Failed to Check Evidence: %v\n", err)
		return
	}

	publicKey := report.Data
	uniqueId := report.UniqueID

	err = CheckSignature(log.Signature, &sign.LogPlain{
		Repository: log.Repository,
		Evidence:   log.Evidence,
		CommitId:   log.CommitID,
	}, monitor.ATPublicKey)

	if err != nil {
		fmt.Printf("Failed to Check Signature: %v\n", err)
		return
	}

	err = CheckSourceHash(log, uniqueId)
	if err != nil {
		fmt.Printf("Failed to Check Source Hash: %v\n", err)
		return
	}

	ta, err := monitor.RegisterTA(publicKey)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	atLog, err := monitor.RegisterATLog(uniqueId, log, ta, true)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Inserted: %v\n", atLog)

}

func (monitor *Monitor) MonitorCTLog(entry crtsh.CertificateEntry) {
	exist := monitor.DB.Client.CTLog.Query().
		Where(ctlog.MonitorLogIDEQ(entry.ID)).
		ExistX(*monitor.DB.Ctx)

	if exist {
		return
	}

	unmarshaledPublicKey, isRSA := entry.Certificate.PublicKey.(*rsa.PublicKey)
	if !isRSA {
		fmt.Printf("Violation: %v\n", errPublicKeyIsNotRSA)
		monitor.RevokeIncompletedCTLog(entry.ID, nil)
		return
	}

	publicKeyBuf := x509.MarshalPKCS1PublicKey(unmarshaledPublicKey)

	ta, exist, err := monitor.SelectOrRegisterTA(publicKeyBuf)
	if err != nil {
		fmt.Printf("Violation: %v\n", err)
		monitor.RevokeIncompletedCTLog(entry.ID, ta)
		return
	}

	if !exist {
		fmt.Printf("TA is not found: %x\n", publicKeyBuf)
		monitor.Revoke(ta)
		return
	}

	ctLog, err := monitor.RegisterCTLog(entry.ID, ta)
	if err != nil {
		panic(err)
	}

	err = NotifyUpdate(monitor)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Inserted: %v\n", ctLog)
}
