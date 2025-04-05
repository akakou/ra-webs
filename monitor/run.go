package monitor

import (
	"fmt"

	"github.com/akakou/ra_webs/monitor/ent"
	"github.com/akakou/ra_webs/monitor/ent/taserver"
)

func (monitor *Monitor) Run(pk publicKey, id int) {
	serv, err := monitor.DB.Client.TAServer.
		Query().
		Order(ent.Desc(taserver.FieldID)).
		First(*monitor.DB.Ctx)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	serv.Update().SetMonitorLogID(id).Save(*monitor.DB.Ctx)

	err = monitor.Check(pk, serv)
	if err != nil {
		fmt.Printf("Violation: %v\n", err)
		monitor.Revoke(serv)
		return
	}

	if !serv.IsActive {
		serv.Update().SetIsActive(true).Save(*monitor.DB.Ctx)
	}
}
