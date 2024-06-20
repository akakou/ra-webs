package notify

import (
	"net/http"

	goutils "github.com/akakou/go-utils"
	rawebscore "github.com/akakou/ra_webs/core"
	"github.com/akakou/ra_webs/ttp/core"
	"github.com/labstack/echo/v4"
)

var postSubscribe = goutils.EchoRoute[core.TTP]{
	Method: goutils.POST,
	Path:   rawebscore.API_ROOT + "/subscription",
	F: func(ttp *core.TTP) goutils.EchoRouteFunc {
		return func(c echo.Context) error {
			var data struct {
				Endpoint       string `json:"endpoint"`
				ExpirationTime int    `json:"expirationTime"`
				Keys           struct {
					Auth   string `json:"auth"`
					P256dh string `json:"p256dh"`
				} `json:"keys"`
			}

			if err := c.Bind(&data); err != nil {
				return err
			}

			subscription, err := ttp.DB.Client.Subscription.
				Create().
				SetEndpoint(data.Endpoint).
				SetAuth(data.Keys.Auth).
				SetP256dh(data.Keys.P256dh).
				Save(*ttp.DB.Ctx)

			if err != nil {
				return err
			}

			return c.JSON(http.StatusOK, subscription)
		}
	},
}
