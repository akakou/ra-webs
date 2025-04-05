package monitor

import (
	"bytes"
	"crypto/rsa"
	"fmt"

	"github.com/akakou/ra-webs/monitor/ent"
	"github.com/google/certificate-transparency-go/x509"
)

const TA_SERVER_NOT_FOUND = "ent: ta_server not found"

type publicKey interface{}

func (monitor *Monitor) Check(pk publicKey, serv *ent.TAServer) error {
	unmarshaledPublicKey, isRSA := pk.(*rsa.PublicKey)
	if !isRSA {
		return errPublicKeyNotRSA
	}

	publicKey := x509.MarshalPKCS1PublicKey(unmarshaledPublicKey)
	fmt.Printf("compireing public key:\n%v\n!=%v\n\n", serv.PublicKey, publicKey)

	if !bytes.Equal(serv.PublicKey, publicKey) {
		return fmt.Errorf("%v: %v != %v", errPublicKeyNotMatch, serv.PublicKey, publicKey)
	}

	return nil
}
