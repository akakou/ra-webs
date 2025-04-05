package monitor

import (
	"github.com/akakou/ra-webs/log/api/interfacestruct"
	"github.com/akakou/ra-webs/monitor/ent"
)

func (monitor *Monitor) RegisterCTLog(evidence string, publicKey []byte) (*ent.CTLog, error) {
	taServerCreate := monitor.DB.Client.CTLog.
		Create().
		SetPublicKey(publicKey).
		SetMonitorLogID(0).
		SetIsActive(true)

	server, err := taServerCreate.Save(*monitor.DB.Ctx)
	if err != nil {
		return nil, err
	}

	return server, nil
}

func (monitor *Monitor) RegisterATLog(uniqueId []byte, log *interfacestruct.TA, ctLog *ent.CTLog) (*ent.ATLog, error) {
	codeCreate := monitor.DB.Client.ATLog.
		Create().
		SetEvidence(log.Evidence).
		SetRepository(log.Repository).
		SetCommitID(log.CommitID).
		SetUniqueID(uniqueId).
		SetCtLog(ctLog).
		SetIsActive(true)

	code, err := codeCreate.Save(*monitor.DB.Ctx)

	if err != nil {
		return nil, err
	}

	return code, nil
}
