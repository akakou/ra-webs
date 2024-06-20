package notify

import (
	"fmt"
	"os"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/akakou/ra_webs/ttp/core"
	"github.com/akakou/ra_webs/ttp/ent"
)

const TTL_MAX = 2419200
const DEFAULT_SUBSCRIBER = "ra-webs@example.com"

type BrowserNotify struct {
	VapidPublicKey, VapidPrivateKey string
	Subscriber                      string
	TTL                             int
}

func NewBrowserNotify(vapidPublicKey, vapidPrivateKey, Subscriber string, TTL int) *BrowserNotify {
	return &BrowserNotify{
		VapidPublicKey:  vapidPublicKey,
		VapidPrivateKey: vapidPrivateKey,
		Subscriber:      Subscriber,
		TTL:             TTL,
	}
}

func DefaultBrowserNotify() (*BrowserNotify, error) {
	var err error

	publicKey := os.Getenv("RA_WEBS_VAPID_PUBLIC_KEY")
	privateKey := os.Getenv("RA_WEBS_VAPID_PRIVATE_KEY")

	if publicKey != "" || privateKey != "" {
		privateKey, publicKey, err = webpush.GenerateVAPIDKeys()
	}

	return NewBrowserNotify(privateKey, publicKey, DEFAULT_SUBSCRIBER, TTL_MAX), err
}

func (notify *BrowserNotify) Notify(msg []byte, serv *ent.TAServer, ttp *core.TTP[BrowserNotify]) error {
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
