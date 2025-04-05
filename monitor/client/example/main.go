package main

import (
	"context"

	"github.com/akakou/ra-webs/monitor"
	"github.com/akakou/ra-webs/monitor/ct/crtsh"
	browsernotifier "github.com/akakou/ra-webs/monitor/notifier/browser"
)

func main() {
	ct, err := crtsh.Default(context.Background())

	if err != nil {
		panic(err)
	}

	notifier, err := browsernotifier.Default()
	if err != nil {
		panic(err)
	}

	m, err := monitor.Default(ct, notifier)
	if err != nil {
		panic(err)
	}

	defer m.Close()

	m.Run()
}
