package notify

import (
	"fmt"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/akakou/ra_webs/ttp/core"
	"github.com/akakou/ra_webs/ttp/ent"
	"github.com/labstack/echo/v4"
)

const VIOLATION_MESSAGE = "A violation has been detected at "
const UPDATE_MESSAGE = "A new server has been added at "

func NotifyViolation(serv *ent.TAServer, ttp *core.TTP) error {
	msg := fmt.Sprintf("%s %v", VIOLATION_MESSAGE, serv.Domain)
	return ttp.Notify.Notify([]byte(msg), serv, ttp)
}

func NotifyUpdate(serv *ent.TAServer, ttp *core.TTP) error {
	msg := fmt.Sprintf("%s %v", UPDATE_MESSAGE, serv.Domain)
	return ttp.Notify.Notify([]byte(msg), serv, ttp)
}

func (notify *BrowserNotify) Notify(msg []byte, serv *ent.TAServer, ttp *core.TTP) error {
	subscriptions, err := serv.QuerySubscription().All(*ttp.DB.Ctx)
	if err != nil {
		panic(err)
	}

	return notify.notifyAll(msg, subscriptions)
}

func (notify *BrowserNotify) notifyAll(msg []byte, subscription []*ent.Subscription) error {
	for _, sub := range subscription {
		err := notify.notifyOne(msg, sub)
		if err != nil {
			return err
		}
	}

	return nil
}

func (notify *BrowserNotify) notifyOne(msg []byte, subscription *ent.Subscription) error {
	s := webpush.Subscription{
		Endpoint: subscription.Endpoint,
		Keys: webpush.Keys{
			Auth:   subscription.Auth,
			P256dh: subscription.P256dh,
		},
	}

	resp, err := webpush.SendNotification(msg, &s, &webpush.Options{
		Subscriber:      notify.Subscriber,
		VAPIDPublicKey:  notify.VapidPublicKey,
		VAPIDPrivateKey: notify.VapidPrivateKey,
		TTL:             notify.TTL,
	})

	if err != nil {
		return fmt.Errorf("%v: %v", ERROR_FAILED_TO_NOTIFY, err)
	}

	defer resp.Body.Close()

	return err
}

func (notify *BrowserNotify) Setup(e *echo.Echo, ttp *core.TTP) error {
	postSubscribe.Set(e, ttp)
	getSubscriptionKey(notify).Set(e, ttp)
	return nil
}
