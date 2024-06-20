package notify

import (
	"os"

	"github.com/SherClockHolmes/webpush-go"
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
