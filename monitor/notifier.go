package monitor

import (
	"fmt"
)

const VIOLATION_MESSAGE = "A violation has been detected at "
const UPDATE_MESSAGE = "A update has been registered at "

func NotifyViolation(domain string, monitor *Monitor) error {
	msg := fmt.Sprintf("%s %v", VIOLATION_MESSAGE, domain)
	return monitor.Notifier.Notify([]byte(msg), domain, monitor)
}

func NotifyUpdate(domain string, monitor *Monitor) error {
	msg := fmt.Sprintf("%s %v", UPDATE_MESSAGE, domain)
	return monitor.Notifier.Notify([]byte(msg), domain, monitor)
}

type Notifier interface {
	Notify(msg []byte, domain string, monitor *Monitor) error
}
