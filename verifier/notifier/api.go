package notifier

import (
	"errors"
	"net/http"

	goutils "github.com/akakou/go-utils"
	"github.com/akakou/ra_webs/verifier/core"
	"github.com/akakou/ra_webs/verifier/ent"
	"github.com/akakou/ra_webs/verifier/ent/taserver"
	"github.com/labstack/echo/v4"
)

var postSubscribeApi = goutils.EchoRoute[core.Verifier]{
	Method: goutils.POST,
	Path:   "/subscription",
	F: func(verifier *core.Verifier) goutils.EchoRouteFunc {
		return func(c echo.Context) error {
			var data struct {
				Domain       string `json:"domain"`
				Subscription struct {
					Endpoint       string `json:"endpoint"`
					ExpirationTime int    `json:"expirationTime"`
					Keys           struct {
						Auth   string `json:"auth"`
						P256dh string `json:"p256dh"`
					} `json:"keys"`
				} `json:"subscription"`
			}

			if err := c.Bind(&data); err != nil {
				return err
			}

			tx, err := verifier.DB.Client.Tx(*verifier.DB.Ctx)
			if err != nil {
				return err
			}

			serv, err := tx.TAServer.
				Query().
				Select(taserver.FieldID).
				Where(taserver.DomainEQ(data.Domain)).
				Order(ent.Desc(taserver.FieldID)).
				First(*verifier.DB.Ctx)

			if err != nil {
				rerr := tx.Rollback()
				return errors.Join(err, rerr)
			}

			subscription, err := tx.Subscription.
				Create().
				SetEndpoint(data.Subscription.Endpoint).
				SetAuth(data.Subscription.Keys.Auth).
				SetP256dh(data.Subscription.Keys.P256dh).
				SetServer(serv).
				Save(*verifier.DB.Ctx)

			if err != nil {
				rerr := tx.Rollback()
				return errors.Join(err, rerr)
			}

			err = tx.Commit()

			if err != nil {
				rerr := tx.Rollback()
				return errors.Join(err, rerr)
			}

			return c.JSON(http.StatusOK, subscription)
		}
	},
}

func getSubscriptionConfigApi(notifier *BrowserNotifier) goutils.EchoRoute[core.Verifier] {
	return goutils.EchoRoute[core.Verifier]{
		Method: goutils.GET,
		Path:   "/config/subscription",
		F: func(verifier *core.Verifier) goutils.EchoRouteFunc {
			return func(c echo.Context) error {
				VapidPublicKey := notifier.VapidPublicKey

				var data struct {
					VapidPublicKey string `json:"vapid_public_key"`
				}

				data.VapidPublicKey = VapidPublicKey

				return c.JSON(http.StatusOK, data)
			}
		},
	}
}
