package serv

import (
	"fmt"

	"github.com/akakou/ra_webs/monitor"
	"github.com/akakou/ra_webs/monitor/ent"
	"github.com/labstack/echo/v4"
)

const VIOLATION_MESSAGE = "A violation has been detected at "
const UPDATE_MESSAGE = "A update has been registered at "

func NotifierViolation(domain string, serv *MonitorServer) error {
	msg := fmt.Sprintf("%s %v", VIOLATION_MESSAGE, domain)
	return serv.Notifier.Notify([]byte(msg), domain, serv)
}

func NotifierUpdate(domain string, serv *MonitorServer) error {
	msg := fmt.Sprintf("%s %v", UPDATE_MESSAGE, domain)
	return serv.Notifier.Notify([]byte(msg), domain, serv)
}

type Notifier interface {
	Notify(msg []byte, domain string, monitorServer *MonitorServer) error
	Setup(e *echo.Group, monitor *MonitorServer) error // todo: fix
}

const ERROR_FAILED_TO_NOTIFY = "failed to notify"

func SelectSubscription(domain string, monitor *monitor.Monitor) ([]*ent.Subscription, error) {
	result := []*ent.Subscription{}
	servers, err := monitor.DB.Client.TAServer.Query().All(*monitor.DB.Ctx)
	if err != nil {
		return nil, fmt.Errorf("%v: %v", ERROR_FAILED_TO_NOTIFY, err)
	}

	for _, serv := range servers {
		subscriptions, err := serv.QuerySubscription().All(*monitor.DB.Ctx)
		if err != nil {
			panic(err)
		}

		result = append(result, subscriptions...)
	}

	return result, nil
}
