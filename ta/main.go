package ta

import (
	"crypto/rsa"
	"crypto/tls"
	"fmt"
	"log"
	"strings"

	"github.com/labstack/echo/v4"
)

type RAConfig struct {
	TTPDomain string
	Domain    string
	Email     string
	JSFolder  string
}

type RA struct {
	config       *RAConfig
	privKeyStore *privKeyStore
	certStore    *certStore
}

func NewRA(config *RAConfig) *RA {
	return &RA{
		config:       config,
		privKeyStore: &privKeyStore{},
		certStore:    &certStore{},
	}
}

func (ra *RA) TLSConfig() (*tls.Config, error) {
	var privKey *rsa.PrivateKey
	var cert *tls.Certificate

	var err error

	if hasFileExists() {
		privKey, cert, err = ra.Load()
	} else {
		privKey, cert, err = ra.Provisioning()
	}

	if err != nil {
		return nil, err
	}

	cert.PrivateKey = privKey

	tlsConfig := tls.Config{
		Certificates: []tls.Certificate{
			*cert,
		},
	}

	return &tlsConfig, nil

}

func (ra *RA) Middleware() func(echo.HandlerFunc) echo.HandlerFunc {
	provisioner := scProvisioner{
		privateKey: ra.privKeyStore.privKey,
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			r := c.Request()

			if r.Header.Get("X-TEE") != "1" {
				u := strings.Split(r.URL.Path, "/")

				if len(u) >= 1 && u[1] == RA_WEBS_FOLDER[1:] {
					if err := next(c); err != nil {
						c.Error(err)
					}
					return nil
				} else {
					return fmt.Errorf("Middleware: X-TEE header is not set (path: %v)", r.URL.Path)
				}
			}

			var sc *secureChannel
			var err error

			// c.Response().Before(func() {
			sc, err = decryptMiddleware(c, provisioner)
			if err != nil {
				c.Error(err)
			}
			fmt.Printf("\nresult: %v %v\n", c.Request().Method, c.Request().URL.Path)

			// })

			c.Response().After(func() {
				c.Response().Header().Set("X-TEE", "1")
				// err := encryptMiddlware(c, sc)

				conn, rw, err := c.Response().Hijack()
				if err != nil {
					c.Error(err)
				}

				defer conn.Close()
				// rw.Peek(int(c.Response().Size))

				// r := http.Response{
				// 	StatusCode: c.Response().Status,
				// 	Header:     c.Response().Header(),
				// }
				// err = r.Write(rw)
				// // _, err = rw.Write([]byte("hello"))
				// if err != nil {
				// 	c.Error(err)
				// }

				// err = rw.Flush()
				// if err != nil {
				// 	c.Error(err)
				// }

				fmt.Print("_", sc, rw)
				// if err != nil {
				// fmt.Printf("err: %v\n", err)
				// c.Error(err)
				// }
			})

			if err := next(c); err != nil {
				c.Error(err)
			}

			// err = encryptMiddlware(c, sc)
			// fmt.Printf("err: %v\n", err)
			// if err != nil {
			// 	c.Error(err)
			// }

			log.Println("before action")

			return nil
		}
	}
}

func (ra *RA) EndPoints(e *echo.Echo) {
	fmt.Printf("%v\n", PUBLIC_KEY_ENDPOINT)
	fmt.Printf("%v\n", SW_ENTRY_ENDPOINT)
	fmt.Printf("%v\n", RA_WEBS_FOLDER)

	ra.makeHTMLEndpoint(e)
	ra.makeJSEndpoint("crypto.js", e, false)
	ra.makeJSEndpoint("entry.js", e, false)
	ra.makeJSEndpoint("sw.js", e, true)
	ra.makePubKeyEndpint(e)

	e.Static(STATIC_FOLDER, ra.config.JSFolder)
}
