package notify

import (
	"net/http"

	goutils "github.com/akakou/go-utils"
	rawebscore "github.com/akakou/ra_webs/core"
	"github.com/akakou/ra_webs/ttp/core"
	"github.com/akakou/ra_webs/ttp/ent"
	"github.com/akakou/ra_webs/ttp/ent/taserver"
	"github.com/labstack/echo/v4"
)

var postSubscribeApi = goutils.EchoRoute[core.TTP]{
	Method: goutils.POST,
	Path:   rawebscore.API_ROOT + "/subscription",
	F: func(ttp *core.TTP) goutils.EchoRouteFunc {
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

			serv, err := ttp.DB.Client.TAServer.
				Query().
				Where(taserver.DomainEQ(data.Domain)).
				Order(ent.Desc(taserver.FieldID)).
				First(*ttp.DB.Ctx)

			if err != nil {
				return err
			}

			subscription, err := ttp.DB.Client.Subscription.
				Create().
				SetEndpoint(data.Subscription.Endpoint).
				SetAuth(data.Subscription.Keys.Auth).
				SetP256dh(data.Subscription.Keys.P256dh).
				SetServer(serv).
				Save(*ttp.DB.Ctx)

			if err != nil {
				return err
			}

			return c.JSON(http.StatusOK, subscription)
		}
	},
}

func getSubscriptionKeyApi(notify *BrowserNotify) goutils.EchoRoute[core.TTP] {
	return goutils.EchoRoute[core.TTP]{
		Method: goutils.GET,
		Path:   rawebscore.API_ROOT + "/subscription_key",
		F: func(ttp *core.TTP) goutils.EchoRouteFunc {
			return func(c echo.Context) error {
				VapidPublicKey := notify.VapidPublicKey

				var data struct {
					VapidPublicKey string `json:"vapid_public_key"`
				}

				data.VapidPublicKey = VapidPublicKey

				return c.JSON(http.StatusOK, data)
			}
		},
	}
}
