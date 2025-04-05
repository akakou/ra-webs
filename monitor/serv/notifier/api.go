package notifier

import (
	"net/http"

	goutils "github.com/akakou/go-utils"
	"github.com/akakou/ra_webs/monitor/ent"
	"github.com/akakou/ra_webs/monitor/ent/taserver"
	"github.com/akakou/ra_webs/monitor/serv"
	"github.com/labstack/echo/v4"
)

var postSubscribeApi = goutils.EchoRoute[serv.MonitorServer]{
	Method: goutils.POST,
	Path:   "/subscription",
	F: func(server *serv.MonitorServer) goutils.EchoRouteFunc {
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

			serv, err := server.Monitor.DB.Client.TAServer.
				Query().
				Order(ent.Desc(taserver.FieldID)).
				First(*server.Monitor.DB.Ctx)

			if err != nil {
				return err
			}

			subscription, err := server.Monitor.DB.Client.Subscription.
				Create().
				SetEndpoint(data.Subscription.Endpoint).
				SetAuth(data.Subscription.Keys.Auth).
				SetP256dh(data.Subscription.Keys.P256dh).
				SetServer(serv).
				Save(*server.Monitor.DB.Ctx)

			if err != nil {
				return err
			}

			return c.JSON(http.StatusOK, subscription)
		}
	},
}

func getSubscriptionConfigApi(notifier *BrowserNotifier) goutils.EchoRoute[serv.MonitorServer] {
	return goutils.EchoRoute[serv.MonitorServer]{
		Method: goutils.GET,
		Path:   "/config/subscription",
		F: func(server *serv.MonitorServer) goutils.EchoRouteFunc {
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
