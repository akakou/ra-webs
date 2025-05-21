package crtsh

import (
	"context"
	"fmt"
	"time"

	"github.com/akakou/crtsh"
	"github.com/akakou/ctstream/direct"
	"github.com/akakou/ra-webs/monitor"
)

type CrtshMonitor struct {
	Last     int
	Interval time.Duration
}

var DefaultInterval = time.Second * 60

const INITIAL_CT_DOMAIN = "example.com"

func New(interval time.Duration, ctx context.Context) (*CrtshMonitor, error) {
	return &CrtshMonitor{
		Interval: interval,
		Last:     0,
	}, nil
}

func Default(ctx context.Context) (*CrtshMonitor, error) {
	return New(DefaultInterval, ctx)
}

func (a *CrtshMonitor) Run(monitor *monitor.Monitor) {
	for {
		time.Sleep(a.Interval)

		entries, err := crtsh.Fetch(monitor.TADomain, "")
		if err == nil {
			fmt.Printf("ct-log: %v\n", entries)
		} else if err.Error() == direct.ERROR_FAILED_TO_NEW {
			continue
		} else {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		monitor.Monitor(entries)
	}
}
