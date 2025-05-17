package monitor

import (
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"reflect"

	"github.com/akakou/crtsh"
	"github.com/akakou/ra-webs/monitor/ent"
	"github.com/akakou/ra-webs/monitor/ent/ctlog"
	"github.com/akakou/ra-webs/monitor/serviceclient"
)

func (monitor *Monitor) Monitor(ctLogs []crtsh.CertificateEntry) {
	for _, entry := range ctLogs {
		monitor.MonitorOne(entry)
	}
}

func (monitor *Monitor) MonitorOne(ctLog crtsh.CertificateEntry) {
	taLog := monitor.MonitorCTLog(ctLog)
	if taLog == nil {
		return
	}

	taEntry, err := monitor.ServiceClient.Fetch(*taLog.PublicKey)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	monitor.MonitorEvidence(taEntry, taLog)
}

func (monitor *Monitor) MonitorEvidence(taEntry *serviceclient.EvidenceEntry, taLog *ent.TA) {
	var err error

	report, err := CheckEvidence(taEntry.Evidence)
	if err != nil {
		fmt.Printf("Failed to Check Evidence: %v\n", err)
		return
	}

	fmt.Printf("left: %x\n right: \n%x\n\n", report.Data, *taLog.PublicKey)

	if !reflect.DeepEqual(report.Data, *taLog.PublicKey) {
		fmt.Printf("Failed to Check Public Key: %v\n", err)
		return
	}

	uniqueId := report.UniqueID

	fmt.Printf("commit id: %v\n", taEntry.CommitID)

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

	ta, err := monitor.RegisterTA(publicKeyBuf)
	if err != nil {
		panic(err)
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
