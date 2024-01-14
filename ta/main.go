package ta

import (
	"crypto/rsa"
	"crypto/tls"
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
)

type RAConfig struct {
	TTPDomain string
	Domain    string
	Email     string
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

			cipher, err := extractCipher(r)
			if err != nil {
				return fmt.Errorf("aa ")
			}

			sc, err := provisioner.provision(cipher.Key)
			if err != nil {
				return fmt.Errorf("aa ")
			}

			plain, err := sc.decrypt(cipher)
			if err != nil {
				return fmt.Errorf("aa ")
			}

			req, err := reqFromJson(plain, r)
			if err != nil {
				return fmt.Errorf("aa ")
			}

			c.SetRequest(req)

			log.Println("before action")
			if err := next(c); err != nil {
				c.Error(err)
			}
			log.Println("after action")

			return nil
		}
	}
}
