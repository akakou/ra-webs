package ttp

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"fmt"
	"math/big"
	"os"

	"time"

	"github.com/akakou/ra_webs/core"
)

var config_path = "ca.json"

type ca struct {
	PrivateKey  *rsa.PrivateKey
	Certificate *x509.Certificate
}

type caConfig struct {
	PrivateKey  []byte `json:"private_key"`
	Certificate []byte `json:"certificate"`
}

func newCA(privateKey *rsa.PrivateKey, certificate *x509.Certificate) *ca {
	return &ca{
		PrivateKey:  privateKey,
		Certificate: certificate,
	}
}

func initCA(subject pkix.Name) (*ca, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	template := x509.Certificate{
		SerialNumber:          big.NewInt(2019),
		Subject:               subject,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)

	if err != nil {
		return nil, err
	}

	certificate, err := x509.ParseCertificate(certBytes)
	if err != nil {
		return nil, err
	}

	return newCA(privateKey, certificate), nil
}

func saveCA(ca *ca) error {
	raw := ca.Certificate.Raw

	config := &caConfig{
		PrivateKey:  x509.MarshalPKCS1PrivateKey(ca.PrivateKey),
		Certificate: raw,
	}

	file, err := os.OpenFile(config_path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	buf, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("error marshalling config: %v", err)
	}

	_, err = file.Write(buf)
	if err != nil {
		fmt.Println("error writing to file:", err)
		return err
	}

	return nil
}

func loadCA() (*ca, error) {
	file, err := os.Open(config_path)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	config := &caConfig{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(config)
	if err != nil {
		return nil, fmt.Errorf("error decoding config: %v", err)
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(config.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("error parsing private key: %v", err)
	}

	certificate, err := x509.ParseCertificate(config.Certificate)
	if err != nil {
		return nil, fmt.Errorf("error parsing certificate: %v", err)
	}

	return newCA(privateKey, certificate), nil
}

func (ca *ca) sign(quote string, subject pkix.Name, publicKey *rsa.PublicKey) ([]byte, error) {
	template := x509.Certificate{
		SerialNumber: big.NewInt(2020),
		Subject:      subject,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(0, 1, 0),
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		ExtraExtensions: []pkix.Extension{
			{
				Id:       core.X509_EXTENSION_LABEL,
				Critical: false,
				Value:    []byte(quote),
			},
		},
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, &template, ca.Certificate, publicKey, ca.PrivateKey)
	if err != nil {
		return nil, err
	}

	return certBytes, nil
}
