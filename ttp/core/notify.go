package core

import (
	"fmt"

	"github.com/akakou/ra_webs/ttp/ent"
	"github.com/akakou/ra_webs/ttp/ent/taserver"
	"github.com/labstack/echo/v4"
)

const VIOLATION_MESSAGE = "A violation has been detected at "
const UPDATE_MESSAGE = "A new server has been added at "

func NotifyViolation(domain string, ttp *TTP) error {
	msg := fmt.Sprintf("%s %v", VIOLATION_MESSAGE, domain)
	return ttp.Notify.Notify([]byte(msg), domain, ttp)
}

func NotifyUpdate(domain string, ttp *TTP) error {
	msg := fmt.Sprintf("%s %v", UPDATE_MESSAGE, domain)
	return ttp.Notify.Notify([]byte(msg), domain, ttp)
}

type Notify interface {
	Notify(msg []byte, domain string, ttp *TTP) error
	Setup(e *echo.Echo, ttp *TTP) error // todo: fix
}

const ERROR_FAILED_TO_NOTIFY = "failed to notify"

func SelectSubscription(domain string, ttp *TTP) ([]*ent.Subscription, error) {
	result := []*ent.Subscription{}
	servers, err := ttp.DB.Client.TAServer.Query().Where(taserver.DomainEQ(domain)).All(*ttp.DB.Ctx)
	if err != nil {
		return nil, fmt.Errorf("%v: %v", ERROR_FAILED_TO_NOTIFY, err)
	}

	for _, serv := range servers {
		subscriptions, err := serv.QuerySubscription().All(*ttp.DB.Ctx)
		if err != nil {
			panic(err)
		}

		result = append(result, subscriptions...)
	}

	return result, nil
}
