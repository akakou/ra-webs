package monitor

import (
	"fmt"

	"github.com/akakou/ra-webs/log/api/io"
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

func (monitor *Monitor) RevokeIncompletedCTLog(ctLogId int, ta *ent.TA) {
	var err error
	if ta == nil {
		ta, err = monitor.RegisterTA([]byte("no public key"))
		if err != nil {
			panic(err)
		}
	}

	fmt.Printf("ta: %v", ta)
	_, err = monitor.RegisterCTLog(ctLogId, ta)
	if err != nil {
		panic(err)
	}

	monitor.Revoke(ta)
}

func (monitor *Monitor) RegisterIncompletedATLog(log *io.TA) {
	ta, err := monitor.RegisterTA([]byte(""))
	if err != nil {
		panic(err)
	}

	_, err = monitor.RegisterATLog([]byte{}, log, ta, false)
	if err != nil {
		panic(err)
	}
}

func (monitor *Monitor) RegisterBrokenATLog(log *io.TA) {
	ta, err := monitor.RegisterTA([]byte(""))
	if err != nil {
		panic(err)
	}

	_, err = monitor.RegisterATLog([]byte{}, log, ta, false)
	if err != nil {
		panic(err)
	}
}
