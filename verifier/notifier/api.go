package notifier

import (
	"net/http"

	goutils "github.com/akakou/go-utils"
	rawebscore "github.com/akakou/ra_webs/core"
	"github.com/akakou/ra_webs/verifier/core"
	"github.com/akakou/ra_webs/verifier/ent"
	"github.com/akakou/ra_webs/verifier/ent/taserver"
	"github.com/labstack/echo/v4"
)

var postSubscribeApi = goutils.EchoRoute[core.Verifier]{
	Method: goutils.POST,
	Path:   rawebscore.API_ROOT + "/subscription",
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

			serv, err := verifier.DB.Client.TAServer.
				Query().
				Where(taserver.DomainEQ(data.Domain)).
				Order(ent.Desc(taserver.FieldID)).
				First(*verifier.DB.Ctx)

			if err != nil {
				return err
			}

			subscription, err := verifier.DB.Client.Subscription.
				Create().
				SetEndpoint(data.Subscription.Endpoint).
				SetAuth(data.Subscription.Keys.Auth).
				SetP256dh(data.Subscription.Keys.P256dh).
				SetServer(serv).
				Save(*verifier.DB.Ctx)

			if err != nil {
				return err
			}

			return c.JSON(http.StatusOK, subscription)
		}
	},
}

func getSubscriptionKeyApi(notifier *BrowserNotifier) goutils.EchoRoute[core.Verifier] {
	return goutils.EchoRoute[core.Verifier]{
		Method: goutils.GET,
		Path:   rawebscore.API_ROOT + "/subscription_key",
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
