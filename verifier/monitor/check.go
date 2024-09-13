package monitor

import (
	"bytes"
	"crypto/rsa"
	"fmt"

	"github.com/akakou/ra_webs/verifier/ent"
	"github.com/google/certificate-transparency-go/x509"
)

const TA_SERVER_NOT_FOUND = "ent: ta_server not found"

func Check(pk interface{}, serv *ent.TAServer) error {
	unmarshaledPublicKey, isRSA := pk.(*rsa.PublicKey)
	if !isRSA {
		return fmt.Errorf("%s", ERROR_PUBLIC_KEY_NOT_RSA)
	}

	publicKey := x509.MarshalPKCS1PublicKey(unmarshaledPublicKey)
	fmt.Printf("compireing public key:\n%v\n!=%v\n\n", serv.PublicKey, publicKey)

	if !bytes.Equal(serv.PublicKey, publicKey) {
		return fmt.Errorf("%v: %v != %v", ERROR_PUBLIC_KEY_NOT_MATCH, serv.PublicKey, publicKey)
	}

	return nil
}
