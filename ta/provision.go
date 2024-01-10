package ta

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/akakou/ra_webs/core"
)

var SCHEME = "https://"

func (ra *RA) Provisioning(pks *privKeyStore, cs *certStore) (*rsa.PrivateKey, *tls.Certificate, error) {
	privKey, pubKey, cert, err := ra.generateKeyPair()
	if err != nil {
		return nil, nil, err
	}

	// attestation, err := raConfig.attest(tlsConfig)
	// if err != nil {
	// 	return nil, err
	// }
	attestation := "attestation"

	pubKeyBytes, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return nil, nil, err
	}

	err = ra.registerToTTP(pubKeyBytes, attestation)
	if err != nil {
		return nil, nil, err
	}

	pks.Set(privKey)
	pks.Store()

	cs.Set(cert)
	cs.Store()

	return privKey, cert, nil
}

func (ra *RA) generateKeyPair() (*rsa.PrivateKey, *rsa.PublicKey, *tls.Certificate, error) {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)

	pubKey := &privKey.PublicKey

	if err != nil {
		return nil, nil, nil, err
	}

	template := &x509.Certificate{
		SerialNumber: &big.Int{},
		Subject:      pkix.Name{CommonName: ra.config.Domain},
		NotAfter:     time.Now().Add(time.Hour),
		DNSNames:     []string{ra.config.Domain},
	}

	rawCert, err := x509.CreateCertificate(rand.Reader, template, template, pubKey, privKey)

	if err != nil {
		return nil, nil, nil, err
	}

	cert := tls.Certificate{
		Certificate: [][]byte{rawCert},
		PrivateKey:  privKey,
	}

	return privKey, pubKey, &cert, nil
}

func (ra *RA) registerToTTP(publicKey []byte, attestation string) error {
	publicKeyHashBytes := sha256.Sum256(publicKey)
	publicKeyHash := hex.EncodeToString(publicKeyHashBytes[:])

	provReq := core.TAInfo{
		Attestation:   attestation,
		PublicKeyHash: publicKeyHash,
		Domain:        ra.config.Domain,
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
