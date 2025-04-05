package notifier

import (
	"fmt"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/akakou/ra-webs/monitor/ent"
	"github.com/akakou/ra-webs/monitor/serv"
	"github.com/labstack/echo/v4"
)

func (notifier *BrowserNotifier) Notify(msg []byte, domain string, server *serv.MonitorServer) error {
	subscriptions, err := serv.SelectSubscription(domain, server.Monitor)
	if err != nil {
		return err
	}

	err = notifier.notifyAll(msg, subscriptions)
	return err
}

func (notifier *BrowserNotifier) notifyAll(msg []byte, subscription []*ent.Subscription) error {
	for _, sub := range subscription {
		err := notifier.notifyOne(msg, sub)
		if err != nil {
			return err
		}
	}

	return nil
}

func (notifier *BrowserNotifier) notifyOne(msg []byte, subscription *ent.Subscription) error {
	s := webpush.Subscription{
		Endpoint: subscription.Endpoint,
		Keys: webpush.Keys{
			Auth:   subscription.Auth,
			P256dh: subscription.P256dh,
		},
	}

	resp, err := webpush.SendNotification(msg, &s, &webpush.Options{
		Subscriber:      notifier.Subscriber,
		VAPIDPublicKey:  notifier.VapidPublicKey,
		VAPIDPrivateKey: notifier.VapidPrivateKey,
		TTL:             notifier.TTL,
	})

	if err != nil {
		return fmt.Errorf("%v: %v", ERROR_FAILED_TO_NOTIFY, err)
	}

	defer resp.Body.Close()

	return err
}

func (notifier *BrowserNotifier) Setup(e *echo.Group, server *serv.MonitorServer) error {
	postSubscribeApi.Set(e, server)
	getSubscriptionConfigApi(notifier).Set(e, server)
	return nil
}
