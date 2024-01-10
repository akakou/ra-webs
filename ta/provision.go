package ta

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/akakou/ra_webs/core"
)

var SCHEME = "https://"
var KEY_SIZE = 2048
var USE_ACME = true

func (ra *RA) Provisioning() (*rsa.PrivateKey, *tls.Certificate, error) {
	privKey, pubKey, err := ra.generateKeyPair()
	if err != nil {
		return nil, nil, err
	}

	var cert *tls.Certificate

	acmeConfig := ACMEConfig{
		Email:      ra.config.Email,
		Domain:     ra.config.Domain,
		PrivateKey: privKey,
	}

	acmeContext, err := initACMEContext(&acmeConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("Provisioning: initACMEContext: %v", err)
	}

	cert, err = acmeContext.issueCert()
	if err != nil {
		return nil, nil, fmt.Errorf("Provisioning: issueCert: %v", err)
	}

	attestation, err := attestateByAzure(pubKey)
	if err != nil {
		return nil, nil, fmt.Errorf("Provisioning: attestateByAzure: %v", err)
	}

	pubKeyBytes, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return nil, nil, err
	}

	err = ra.registerToTTP(pubKeyBytes, attestation)
	if err != nil {
		return nil, nil, err
	}

	ra.privKeyStore.privKey = privKey
	ra.privKeyStore.Store()

	ra.certStore.cert = cert
	ra.certStore.Store()

	return privKey, cert, nil
}

func (ra *RA) generateKeyPair() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privKey, err := rsa.GenerateKey(rand.Reader, KEY_SIZE)
	if err != nil {
		return nil, nil, err
	}

	pubKey := &privKey.PublicKey

	return privKey, pubKey, nil
}

func (ra *RA) registerToTTP(publicKey []byte, attestation string) error {
	provReq := core.ProvisionRequest{
		Attestation: attestation,
		Domain:      ra.config.Domain,
	}

	body, _ := json.Marshal(provReq)
	buf := bytes.NewBuffer(body)

	req, err := http.NewRequest("POST", SCHEME+ra.config.TTPDomain+"/provision", buf)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return nil
	} else {
		return fmt.Errorf("TTP returned %s", resp.Status)
	}
}
