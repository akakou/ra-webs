package monitor

import (
	"github.com/akakou/ra_webs/monitor"
	"github.com/akakou/ra_webs/monitor/ent"
)

func revoke(serv *ent.TAServer, monitor *monitor.Monitor) {
	// err := core.NotifierViolation(monitor.Domain, monitor)
	// if err != nil {
	// 	panic(err)
	// }

	monitor.DB.Client.TAViolation.Create().
		SetServer(serv).
		SaveX(*monitor.DB.Ctx)

	service := serv.QueryService().OnlyX(*monitor.DB.Ctx)

	service.Update().SetIsActive(false).SaveX(*monitor.DB.Ctx)
}
