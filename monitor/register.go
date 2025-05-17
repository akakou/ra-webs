package monitor

import (
	"github.com/akakou/ra-webs/monitor/ent"
	"github.com/akakou/ra-webs/monitor/serviceclient"
)

func (monitor *Monitor) RegisterTA(publicKey []byte) *ent.TA {
	ta := monitor.DB.Client.TA.
		Create().
		SetPublicKey(publicKey).
		SaveX(*monitor.DB.Ctx)

	return ta
}

func (monitor *Monitor) RegisterCTLog(ctLogId int, ta *ent.TA, isActive bool) *ent.CTLog {
	ctLog := monitor.DB.Client.CTLog.
		Create().
		SetMonitorLogID(ctLogId).
		SetTa(ta).
		SaveX(*monitor.DB.Ctx)

	return ctLog
}

func (monitor *Monitor) RegisterATLog(uniqueId []byte, entry *serviceclient.EvidenceEntry, ta *ent.TA, active bool) *ent.ATLog {
	atLog := monitor.DB.Client.ATLog.
		Create().
		SetEvidence(entry.Evidence).
		SetRepository(entry.Repository).
		SetCommitID(entry.CommitID).
		SetUniqueID(uniqueId).
		SetTa(ta).
		SaveX(*monitor.DB.Ctx)

	return atLog
}
