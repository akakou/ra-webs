package ta

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
)

var STATIC_FOLDER = "./static"

func (ra *RA) MakeServiceWorker() (string, error) {
	template := "const PUBLIC_KEY = '%v';"

	pubKey := ra.privKeyStore.privKey.Public().(*rsa.PublicKey)
	rawPubKey, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return "", err
	}

	base64PubKey := base64.StdEncoding.EncodeToString(rawPubKey)

	js := fmt.Sprintf(template, base64PubKey)

	return js, nil
}
