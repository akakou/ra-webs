package monitor

import (
	"fmt"
)

const VIOLATION_MESSAGE = "A violation has been detected at "
const UPDATE_MESSAGE = "A update has been registered at "

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func NotifyViolationX(monitor *Monitor) {
	panicIfErr(NotifyViolation(monitor))
}

func NotifyUpdateX(monitor *Monitor) {
	panicIfErr(NotifyUpdate(monitor))
}

func NotifyViolation(monitor *Monitor) error {
	msg := fmt.Sprintf("%s %v", VIOLATION_MESSAGE, monitor.TADomain)
	return monitor.Notifier.Notify([]byte(msg), monitor)
}

func NotifyUpdate(monitor *Monitor) error {
	msg := fmt.Sprintf("%s %v", UPDATE_MESSAGE, monitor.TADomain)
	return monitor.Notifier.Notify([]byte(msg), monitor)
}

type Notifier interface {
	Notify(msg []byte, monitor *Monitor) error
}
