package ta

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

var RA_WEBS_FOLDER = "/ra-webs"
var JS_FOLDER = RA_WEBS_FOLDER + "/static"

var PUBLIC_KEY_ENDPOINT = RA_WEBS_FOLDER + "/public_key.js"
var SW_ENTRY_ENDPOINT = RA_WEBS_FOLDER + "/entry"
var STATIC_FOLDER = RA_WEBS_FOLDER + "/static"
var SW_ENTRY_ENDPOINT_JS = STATIC_FOLDER + "/entry.js"

func readFile(name string) string {
	file, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	bytes, _ := io.ReadAll(file)
	html := string(bytes)

	return html
}

func (ra *RA) makeHTMLEndpoint(e *echo.Echo) {
	e.GET(SW_ENTRY_ENDPOINT, func(c echo.Context) error {
		template := fmt.Sprintf("<script src='%v'></script>", SW_ENTRY_ENDPOINT_JS)
		return c.HTML(http.StatusOK, template)
	})
}

func (ra *RA) makeJSEndpoint(name string, e *echo.Echo, isSW bool) {
	endPointPath := STATIC_FOLDER + "/" + name
	fsPath := ra.config.JSFolder + "/" + name

	e.GET(endPointPath, func(c echo.Context) error {
		js := readFile(fsPath)
		headers := c.Response().Header()

		if isSW {
			headers.Set("Service-Worker-Allowed", "/")

		}
		headers.Set("Content-Type", "application/javascript")

		return c.String(http.StatusOK, js)
	})
}

func (ra *RA) makePubKeyEndpint(e *echo.Echo) {
	e.GET(PUBLIC_KEY_ENDPOINT, func(c echo.Context) error {
		template := "const PUBLIC_KEY = '%v'"

		pubKey := ra.privKeyStore.privKey.Public().(*rsa.PublicKey)
		rawPubKey, err := x509.MarshalPKIXPublicKey(pubKey)
		if err != nil {
			return err
		}

		base64PubKey := base64.StdEncoding.EncodeToString(rawPubKey)

		js := fmt.Sprintf(template, base64PubKey)
		c.Response().Header().Set("Content-Type", "application/javascript")
		return c.String(http.StatusOK, js)
	})
}
