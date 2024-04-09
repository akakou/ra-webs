package ta

import (
	"crypto/rsa"
	"crypto/x509"

	"github.com/akakou/ra_webs/core"
)

func attestPublicKey(ap *TA) (string, error) {
	publicKey := ap.PrivateKey.Public()
	publicKeyBuf := x509.MarshalPKCS1PublicKey(publicKey.(*rsa.PublicKey))

	quote, err := core.AttestByAzure(publicKeyBuf)
	return quote, err
}
