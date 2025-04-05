package monitor

import (
	"github.com/akakou/ra-webs/log/api/io"
	"github.com/akakou/ra-webs/monitor/ent"
	"github.com/akakou/ra-webs/monitor/ent/ta"
)

func (monitor *Monitor) SelectOrRegisterTA(publicKey []byte) (*ent.TA, bool, error) {
	ta, err := monitor.DB.Client.TA.Query().
		Where(ta.PublicKeyEQ(publicKey)).
		Only(*monitor.DB.Ctx)

	if err != nil {
		return nil, false, err
	}

	if ta != nil {
		return ta, true, nil
	}

	ta, err = monitor.RegisterTA(publicKey)

	if err != nil {
		return nil, false, err
	}

	return ta, false, nil
}

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

func (monitor *Monitor) RegisterCTLog(ctLogId int, ta *ent.TA) (*ent.CTLog, error) {
	ctLogCreate := monitor.DB.Client.CTLog.
		Create().
		SetMonitorLogID(ctLogId).
		SetTa(ta).
		SetIsActive(true)

	ctLog, err := ctLogCreate.Save(*monitor.DB.Ctx)
	if err != nil {
		return nil, err
	}

	return ctLog, nil
}

func (monitor *Monitor) RegisterATLog(uniqueId []byte, log *io.TA, ta *ent.TA) (*ent.ATLog, error) {
	atLogCreate := monitor.DB.Client.ATLog.
		Create().
		SetEvidence(log.Evidence).
		SetRepository(log.Repository).
		SetCommitID(log.CommitID).
		SetUniqueID(uniqueId).
		SetSignature(log.Signature).
		SetTa(ta).
		SetIsActive(true)

	atLog, err := atLogCreate.Save(*monitor.DB.Ctx)

	if err != nil {
		return nil, err
	}

	return atLog, nil
}
