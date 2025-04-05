package monitor

import (
	"github.com/akakou/ra-webs/log/api/interfacestruct"
	"github.com/akakou/ra-webs/monitor/ent"
)

func (monitor *Monitor) RegisterServer(evidence string, publicKey []byte, code *ent.TACode) (*ent.TAServer, error) {
	taServerCreate := monitor.DB.Client.TAServer.
		Create().
		SetCode(code).
		SetPublicKey(publicKey).
		SetQuote(evidence).
		SetMonitorLogID(0).
		SetIsActive(true)

	server, err := taServerCreate.Save(*monitor.DB.Ctx)
	if err != nil {
		return nil, err
	}

	return server, nil
}

func (monitor *Monitor) RegisterCode(uniqueId []byte, log *interfacestruct.TA) (*ent.TACode, error) {
	codeCreate := monitor.DB.Client.TACode.
		Create().
		SetRepository(log.Repository).
		SetCommitID(log.CommitID).
		SetUniqueID(uniqueId).
		SetIsActive(true)

	code, err := codeCreate.Save(*monitor.DB.Ctx)

	if err != nil {
		return nil, err
	}

	return code, nil
}
