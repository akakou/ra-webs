package ta

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"io"
	"os"
	"strings"
)

var STATIC_FOLDER = "./static"

func (ra *RA) MakeServiceWorker() (string, error) {
	file, err := os.Open(STATIC_FOLDER + "/sw_template.js")
	if err != nil {
		return "", err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	pubKey := ra.privKeyStore.privKey.Public().(*rsa.PublicKey)
	rawPubKey := x509.MarshalPKCS1PublicKey(pubKey)

	base64PubKey := base64.StdEncoding.EncodeToString(rawPubKey)

	template := string(bytes)
	js := strings.Replace(template, "{{PUBLIC_KEY}}", base64PubKey, 1)

	return js, nil
}
