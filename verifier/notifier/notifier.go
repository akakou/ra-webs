package notifier

import (
	"fmt"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/akakou/ra_webs/verifier/core"
	"github.com/akakou/ra_webs/verifier/ent"
	"github.com/labstack/echo/v4"
)

func (notifier *BrowserNotifier) Notifier(msg []byte, domain string, verifier *core.Verifier) error {
	subscriptions, err := core.SelectSubscription(domain, verifier)
	if err != nil {
		return err
	}

	err = notifier.notifierAll(msg, subscriptions)
	return err
}

func (notifier *BrowserNotifier) notifierAll(msg []byte, subscription []*ent.Subscription) error {
	for _, sub := range subscription {
		err := notifier.notifierOne(msg, sub)
		if err != nil {
			return err
		}
	}

	return nil
}

func (notifier *BrowserNotifier) notifierOne(msg []byte, subscription *ent.Subscription) error {
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

func (notifier *BrowserNotifier) Setup(e *echo.Echo, verifier *core.Verifier) error {
	postSubscribeApi.Set(e, verifier)
	getSubscriptionKeyApi(notifier).Set(e, verifier)
	return nil
}
