package monitor

import (
	"github.com/akakou/ra-webs/log/api/interfacestruct"
	"github.com/akakou/ra-webs/monitor/ent"
)

func (monitor *Monitor) Revoke(ta *ent.TA) {
	err := NotifyViolation(monitor.Domain, monitor)
	if err != nil {
		panic(err)
	}

	monitor.DB.Client.Violation.Create().
		SetTa(ta).
		SaveX(*monitor.DB.Ctx)
}

func (monitor *Monitor) RevokeWithEmpty() {
	ta, err := monitor.RegisterTA([]byte(""))
	if err != nil {
		panic(err)
	}

	monitor.Revoke(ta)
}

func (monitor *Monitor) RevokeIncompletedCTLog(ctLogId int) {
	ta, err := monitor.RegisterTA([]byte(""))
	if err != nil {
		panic(err)
	}

	_, err = monitor.RegisterCTLog(ctLogId, ta)

	monitor.Revoke(ta)
}

func (monitor *Monitor) RevokeIncompletedATLog(log *interfacestruct.TA) {
	ta, err := monitor.RegisterTA([]byte(""))
	if err != nil {
		panic(err)
	}

	_, err = monitor.RegisterATLog([]byte{}, log, ta)

	monitor.Revoke(ta)
}

func (monitor *Monitor) RevokeATLog(log *interfacestruct.TA) {
	ta, err := monitor.RegisterTA([]byte(""))
	if err != nil {
		panic(err)
	}

	_, err = monitor.RegisterATLog([]byte{}, log, ta)

	monitor.Revoke(ta)
}
