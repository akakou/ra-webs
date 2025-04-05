package monitor

import (
	"github.com/akakou/ra-webs/monitor/ent"
)

func (monitor *Monitor) Revoke(serv *ent.CTLog) {
	err := NotifyViolation(monitor.Domain, monitor)
	if err != nil {
		panic(err)
	}

	monitor.DB.Client.TAViolation.Create().
		SetCtLog(serv).
		SaveX(*monitor.DB.Ctx)
}
