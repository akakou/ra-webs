package monitor

import (
	"github.com/akakou/ra-webs/monitor/ent"
	"github.com/akakou/ra-webs/monitor/serviceclient"
)

func (monitor *Monitor) RegisterTA(publicKey []byte) (*ent.TA, error) {
	taCreate := monitor.DB.Client.TA.
		Create().
		SetPublicKey(publicKey)

	ta, err := taCreate.Save(*monitor.DB.Ctx)
	if err != nil {
		return nil, err
	}

	return ta, nil
}

func (monitor *Monitor) RegisterCTLog(ctLogId int, ta *ent.TA, isActive bool) (*ent.CTLog, error) {
	ctLogCreate := monitor.DB.Client.CTLog.
		Create().
		SetMonitorLogID(ctLogId).
		SetTa(ta).
		SetIsActive(isActive)

	ctLog, err := ctLogCreate.Save(*monitor.DB.Ctx)
	if err != nil {
		return nil, err
	}

	return ctLog, nil
}

func (monitor *Monitor) RegisterATLog(uniqueId []byte, entry *serviceclient.EvidenceEntry, ta *ent.TA, active bool) (*ent.ATLog, error) {
	atLogCreate := monitor.DB.Client.ATLog.
		Create().
		SetEvidence(entry.Evidence).
		SetRepository(entry.Repository).
		SetCommitID(entry.CommitID).
		SetUniqueID(uniqueId).
		SetTa(ta)

	atLog, err := atLogCreate.Save(*monitor.DB.Ctx)

	if err != nil {
		return nil, err
	}

	return atLog, nil
}
