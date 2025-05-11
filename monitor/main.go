package monitor

import (
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"reflect"

	"github.com/akakou/crtsh"
	"github.com/akakou/ra-webs/monitor/ent"
	"github.com/akakou/ra-webs/monitor/ent/ctlog"
	"github.com/akakou/ra-webs/service/api/io"
)

func (monitor *Monitor) Monitor(ctLogs []crtsh.CertificateEntry) {
	for _, entry := range ctLogs {
		monitor.MonitorOne(entry)
	}
}

func (monitor *Monitor) MonitorOne(ctLog crtsh.CertificateEntry) {
	taLog := monitor.MonitorCTLog(ctLog)
	taEntry, err := monitor.ServiceClient.Fetch(*taLog.PublicKey)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	monitor.MonitorEvidence(taEntry, taLog)
}

func (monitor *Monitor) MonitorEvidence(taEntry *io.TA, taLog *ent.TA) {
	var err error

	report, err := CheckEvidence(taEntry.Evidence)
	if err != nil {
		fmt.Printf("Failed to Check Evidence: %v\n", err)
		return
	}

	if !reflect.DeepEqual(report.Data, taEntry.PublicKey) {
		fmt.Printf("Failed to Check Public Key: %v\n", err)
		return
	}

	uniqueId := report.UniqueID

	err = CheckSourceHash(taEntry, uniqueId)
	if err != nil {
		fmt.Printf("Failed to Check Source Hash: %v\n", err)
		return
	}

	atLog, err := monitor.RegisterATLog(uniqueId, taEntry, taLog, true)
	if err != nil {
		fmt.Printf("Error (2): %v\n", err)
		return
	}

	fmt.Printf("Inserted: %v\n", atLog)
}

func (monitor *Monitor) MonitorCTLog(entry crtsh.CertificateEntry) *ent.TA {
	exist := monitor.DB.Client.CTLog.Query().
		Where(ctlog.MonitorLogIDEQ(entry.ID)).
		ExistX(*monitor.DB.Ctx)

	if exist {
		return nil
	}

	unmarshaledPublicKey, isRSA := entry.Certificate.PublicKey.(*rsa.PublicKey)
	if !isRSA {
		fmt.Printf("Violation (a): %v\n", errPublicKeyIsNotRSA)
		monitor.RevokeIncompletedCTLog(entry.ID, nil)
		return nil
	}

	publicKeyBuf := x509.MarshalPKCS1PublicKey(unmarshaledPublicKey)

	ta, exist, err := monitor.SelectOrRegisterTA(publicKeyBuf)
	if err != nil {
		fmt.Printf("Violation (b): %v\n", err)
		monitor.RevokeIncompletedCTLog(entry.ID, ta)
		return nil
	}

	if !exist {
		fmt.Printf("TA is not found: %x\n", publicKeyBuf)
		monitor.Revoke(ta)
		return nil
	}

	_, err = monitor.RegisterCTLog(entry.ID, ta, true)
	if err != nil {
		panic(err)
	}

	err = NotifyUpdate(monitor)
	if err != nil {
		panic(err)
	}

	return ta
}
