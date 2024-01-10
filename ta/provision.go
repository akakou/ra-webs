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

func (ra *RA) Provisioning() (*tls.Config, error) {
	tlsConfig, publicKey, err := ra.generateKeyPair()
	if err != nil {
		return nil, err
	}

	// attestation, err := raConfig.attest(tlsConfig)
	// if err != nil {
	// 	return nil, err
	// }
	attestation := "attestation"

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, err
	}

	err = ra.registerToTTP(publicKeyBytes, attestation)
	if err != nil {
		return nil, err
	}

	return tlsConfig, nil
}

func (ra *RA) generateKeyPair() (*tls.Config, *rsa.PublicKey, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		return nil, nil, err
	}

	template := &x509.Certificate{
		SerialNumber: &big.Int{},
		Subject:      pkix.Name{CommonName: ra.config.Domain},
		NotAfter:     time.Now().Add(time.Hour),
		DNSNames:     []string{ra.config.Domain},
	}

	cert, err := x509.CreateCertificate(rand.Reader, template, template, &priv.PublicKey, priv)

	if err != nil {
		return nil, nil, err
	}

	tlsCfg := tls.Config{
		Certificates: []tls.Certificate{
			{
				Certificate: [][]byte{cert},
				PrivateKey:  priv,
			},
		},
	}

	return &tlsCfg, &priv.PublicKey, nil
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
