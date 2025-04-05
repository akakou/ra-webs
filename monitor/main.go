package monitor

import (
	"crypto/rsa"
	"crypto/x509"
	"fmt"

	"github.com/akakou/crtsh"
	"github.com/akakou/ra-webs/log/api/interfacestruct"
	"github.com/akakou/ra-webs/monitor/ent/atlog"
	"github.com/akakou/ra-webs/monitor/ent/ctlog"
)

type publicKey interface{}

func (monitor *Monitor) Monitor(ctLogs []crtsh.CertificateEntry) {
	atLogs, err := monitor.LogClient.Fetch()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		monitor.RevokeIncompletedATLog(&interfacestruct.TA{
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

func (monitor *Monitor) MonitorATLog(log *interfacestruct.TA) {
	var err error

	exist, err := monitor.DB.Client.ATLog.Query().
		Where(atlog.EvidenceEQ(log.Evidence)).
		Exist(*monitor.DB.Ctx)

	if err != nil {
		panic(err)
	}

	if exist {
		return
	}

	report, err := CheckEvidence(log.Evidence)
	if err != nil {
		fmt.Printf("Violation: %v\n", err)
		monitor.RevokeIncompletedATLog(log)
		return
	}

	publicKey := report.Data
	uniqueId := report.Data

	ta, _, err := monitor.SelectOrRegisterTA(publicKey)
	if err != nil {
		fmt.Printf("Violation: %v\n", err)
		monitor.RevokeIncompletedATLog(log)
		return
	}

	_, err = monitor.RegisterATLog(uniqueId, log, ta)
	if err != nil {
		fmt.Printf("Violation: %v\n", err)
		monitor.RevokeIncompletedATLog(log)
		return
	}

	err = CheckSourceHash(log, uniqueId)
	if err != nil {
		fmt.Printf("Violation: %v\n", err)
		monitor.Revoke(ta)
		return
	}
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
