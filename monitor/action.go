package monitor

import (
	"github.com/akakou/ra-webs/monitor/ent"
)

func (monitor *Monitor) Revoke(serv *ent.TAServer) {
	err := NotifyViolation(monitor.Domain, monitor)
	if err != nil {
		panic(err)
	}

	monitor.DB.Client.TAViolation.Create().
		SetServer(serv).
		SaveX(*monitor.DB.Ctx)

	service := serv.QueryService().OnlyX(*monitor.DB.Ctx)

	service.Update().SetIsActive(false).SaveX(*monitor.DB.Ctx)
}

func (monitor *Monitor) Activate(serv *ent.TAServer) {
	err := NotifyUpdate(monitor.Domain, monitor)
	if err != nil {
		panic(err)
	}

	serv.Update().SetIsActive(true).SaveX(*monitor.DB.Ctx)
}
