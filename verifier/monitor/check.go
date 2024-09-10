package monitor

import (
	"bytes"
	"crypto/rsa"
	"fmt"

	"github.com/akakou/ra_webs/verifier/core"
	"github.com/akakou/ra_webs/verifier/ent"
	"github.com/akakou/ra_webs/verifier/ent/taserver"
	"github.com/google/certificate-transparency-go/x509"
)

const TA_SERVER_NOT_FOUND = "ent: ta_server not found"

func check(domain string, pk interface{}, ctLogId string, verifier *core.Verifier) error {
	serv, err := verifier.DB.Client.TAServer.
		Query().
		Where(taserver.DomainEQ(domain)).
		Where(taserver.HasActivated(false)).
		Order(ent.Desc(taserver.FieldID)).
		First(*verifier.DB.Ctx)

	if err == nil {
	} else if err.Error() == TA_SERVER_NOT_FOUND {
		return nil
	} else {
		return fmt.Errorf("%v: %v", ERROR_SELECT_SERVER, err)
	}

	unmarshaledPublicKey, isRSA := pk.(*rsa.PublicKey)
	if !isRSA {
		revoke(serv, verifier)
		return fmt.Errorf("%s", ERROR_PUBLIC_KEY_NOT_RSA)
	}

	publicKey := x509.MarshalPKCS1PublicKey(unmarshaledPublicKey)

	if !bytes.Equal(serv.PublicKey, publicKey) {
		revoke(serv, verifier)
		return fmt.Errorf("%v: %v", ERROR_PUBLIC_KEY_NOT_MATCH, err)
	}

	if !serv.HasActivated {
		UpdateLastLog(serv.Domain, ctLogId, verifier.DB)
		serv.Update().SetHasActivated(true).Save(*verifier.DB.Ctx)
	}

	return nil
}

func Check(cert *x509.Certificate, ctLogId string, verifier *core.Verifier) error {
	for _, domain := range cert.DNSNames {
		err := check(domain, cert.PublicKey, ctLogId, verifier)

		if err != nil {
			return err
		}
	}

	return nil
}
