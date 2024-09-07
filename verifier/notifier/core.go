package notifier

import (
	"os"

	"github.com/SherClockHolmes/webpush-go"
)

const TTL_MAX = 2419200
const DEFAULT_SUBSCRIBER = "ra-webs@example.com"

type BrowserNotifier struct {
	VapidPrivateKey, VapidPublicKey string
	Subscriber                      string
	TTL                             int
}

func NewBrowserNotifier(vapidPrivateKey, vapidPublicKey, Subscriber string, TTL int) *BrowserNotifier {
	return &BrowserNotifier{
		VapidPrivateKey: vapidPrivateKey,
		VapidPublicKey:  vapidPublicKey,
		Subscriber:      Subscriber,
		TTL:             TTL,
	}
}

func DefaultBrowserNotifier() (*BrowserNotifier, error) {
	var err error

	publicKey := os.Getenv("RA_WEBS_VAPID_PUBLIC_KEY")
	privateKey := os.Getenv("RA_WEBS_VAPID_PRIVATE_KEY")

	if publicKey == "" || privateKey == "" {
		privateKey, publicKey, err = webpush.GenerateVAPIDKeys()
	}

	return NewBrowserNotifier(privateKey, publicKey, DEFAULT_SUBSCRIBER, TTL_MAX), err
}
