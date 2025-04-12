package crtsh

import (
	"context"
	"fmt"
	"time"

	"github.com/akakou/crtsh"
	"github.com/akakou/ctstream/direct"
	"github.com/akakou/ra-webs/monitor"
)

var Sleep = time.Second * 10

type CrtshMonitor struct {
	Last     int
	Interval time.Duration
}

var DefaultInterval = 10 * time.Minute

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
		entries, err := crtsh.Fetch(monitor.Domain, "")
		if err == nil {
			fmt.Printf("ct-log: %v\n", entries)
		} else if err.Error() == direct.ERROR_FAILED_TO_NEW {
			return
		} else {
			fmt.Printf("Error: %v\n", err)
			return
		}

		monitor.Monitor(entries)
		time.Sleep(Sleep)
	}
}
