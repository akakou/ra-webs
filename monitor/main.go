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
		monitor.RegisterBrokenATLog(&io.TA{
			Evidence:   "",
			Signature:  []byte(""),
			Repository: "",
			CommitID:   "",
		})

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
		fmt.Printf("Violation: %v\n", err)
		monitor.RegisterIncompletedATLog(log)
		return
	}

	publicKey := report.Data
	uniqueId := report.Data

	ta, _, err := monitor.SelectOrRegisterTA(publicKey)
	if err != nil {
		fmt.Printf("Violation: %v\n", err)
		monitor.RegisterIncompletedATLog(log)
		return
	}

	atLog, err := monitor.RegisterATLog(uniqueId, log, ta, false)
	if err != nil {
		fmt.Printf("Violation: %v\n", err)
		return
	}

	err = CheckSignature(log.Signature, &sign.LogPlain{
		Repository: log.Repository,
		Evidence:   log.Evidence,
		CommitId:   log.CommitID,
	}, monitor.ATPublicKey)

	if err != nil {
		fmt.Printf("Violation: %v\n", err)
		return
	}

	err = CheckSourceHash(log, uniqueId)
	if err != nil {
		fmt.Printf("Violation: %v\n", err)
		return
	}

	atLog.Update().SetIsActive(true).SaveX(*monitor.DB.Ctx)
}

func (monitor *Monitor) MonitorCTLog(entry crtsh.CertificateEntry) {
	exist, err := monitor.DB.Client.CTLog.Query().
		Where(ctlog.MonitorLogIDEQ(entry.ID)).
		Exist(*monitor.DB.Ctx)

	if err != nil {
		panic(err)
	}

	if exist {
		return
	}

	unmarshaledPublicKey, isRSA := entry.Certificate.PublicKey.(*rsa.PublicKey)
	if !isRSA {
		fmt.Printf("Violation: %v\n", err)
		monitor.RevokeIncompletedCTLog(entry.ID)
		return
	}

	publicKeyBuf := x509.MarshalPKCS1PublicKey(unmarshaledPublicKey)

	ta, exist, err := monitor.SelectOrRegisterTA(publicKeyBuf)
	if err != nil {
		fmt.Printf("Violation: %v\n", err)
		monitor.RevokeIncompletedCTLog(entry.ID)
		return
	}

	if exist {
		_, err = monitor.RegisterCTLog(entry.ID, ta)
		if err != nil {
			panic(err)
		}

		err = NotifyUpdate(monitor.Domain, monitor)
		if err != nil {
			panic(err)
		}
	} else {
		monitor.Revoke(ta)
	}
}
