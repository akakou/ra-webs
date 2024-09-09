package core

import (
	"fmt"

	"github.com/akakou/ra_webs/verifier/ent"
	"github.com/akakou/ra_webs/verifier/ent/taserver"
	"github.com/labstack/echo/v4"
)

const VIOLATION_MESSAGE = "A violation has been detected at "
const UPDATE_MESSAGE = "A new server has been added at "

func NotifierViolation(domain string, verifier *Verifier) error {
	msg := fmt.Sprintf("%s %v", VIOLATION_MESSAGE, domain)
	return verifier.Notifier.Notify([]byte(msg), domain, verifier)
}

func NotifierUpdate(domain string, verifier *Verifier) error {
	msg := fmt.Sprintf("%s %v", UPDATE_MESSAGE, domain)
	return verifier.Notifier.Notify([]byte(msg), domain, verifier)
}

type Notifier interface {
	Notify(msg []byte, domain string, verifier *Verifier) error
	Setup(e *echo.Echo, verifier *Verifier) error // todo: fix
}

const ERROR_FAILED_TO_NOTIFY = "failed to notify"

func SelectSubscription(domain string, verifier *Verifier) ([]*ent.Subscription, error) {
	result := []*ent.Subscription{}
	servers, err := verifier.DB.Client.TAServer.Query().Where(taserver.DomainEQ(domain)).All(*verifier.DB.Ctx)
	if err != nil {
		return nil, fmt.Errorf("%v: %v", ERROR_FAILED_TO_NOTIFY, err)
	}

	for _, serv := range servers {
		subscriptions, err := serv.QuerySubscription().All(*verifier.DB.Ctx)
		if err != nil {
			panic(err)
		}

		result = append(result, subscriptions...)
	}

	return result, nil
}
