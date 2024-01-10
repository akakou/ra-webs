package ta

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"log"

	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge/tlsalpn01"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
)

func IssueCertiifcate(user *acmeUser, privateKey *rsa.PrivateKey, domains []string, isProduction bool) (*certificate.Resource, error) {
	var acmeClient *lego.Client
	var err error

	if isProduction {
		acmeClient, err = NewAcmeClient(user)
	} else {
		acmeClient, err = NewStagingAcmeClient(user)
	}

	if err != nil {
		return nil, fmt.Errorf("IssueCertificate: creating ACME Client: %v", err)
	}

	err = SetupProvider(acmeClient)
	if err != nil {
		return nil, fmt.Errorf("IssueCertificate: setting up provider: %v", err)
	}

	err = user.Register(acmeClient)
	if err != nil {
		return nil, fmt.Errorf("IssueCertificate: registering acme user: %v", err)
	}

	cert, err := user.ObtainCertificate(acmeClient, privateKey, domains)
	if err != nil {
		return nil, fmt.Errorf("IssueCertificate: obtaining certificate: %v", err)
	}

	log.Printf("IssueCertificate: obtained certificate for domain: %v\n", cert.Domain)

	return cert, nil
}

// Implhent acme.User interface
type acmeUser struct {
	Email        string
	Registration *registration.Resource
	key          crypto.PrivateKey
}

func (u *acmeUser) GetEmail() string {
	return u.Email
}

func (u acmeUser) GetRegistration() *registration.Resource {
	return u.Registration
}

func (u *acmeUser) GetPrivateKey() crypto.PrivateKey {
	return u.key
}

func NewUser(email string, key crypto.PrivateKey) *acmeUser {
	return &acmeUser{
		Email: email,
		key:   key,
	}
}

func CreateUser(email string) (*acmeUser, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	return NewUser(email, privateKey), nil
}

func NewStagingAcmeClient(user *acmeUser) (*lego.Client, error) {
	config := lego.NewConfig(user)
	config.CADirURL = lego.LEDirectoryStaging

	return lego.NewClient(config)
}

func NewAcmeClient(user *acmeUser) (*lego.Client, error) {
	config := lego.NewConfig(user)
	return lego.NewClient(config)
}

func SetupProvider(client *lego.Client) error {
	return client.Challenge.SetTLSALPN01Provider(tlsalpn01.NewProviderServer("", "443"))
}

func (user *acmeUser) Register(client *lego.Client) error {
	options := registration.RegisterOptions{
		TermsOfServiceAgreed: true,
	}
	reg, err := client.Registration.Register(options)
	if err != nil {
		return err
	}

	user.Registration = reg
	return nil
}

func (user *acmeUser) ObtainCertificate(client *lego.Client, privateKey *rsa.PrivateKey, domains []string) (*certificate.Resource, error) {
	request := certificate.ObtainRequest{
		Domains:    domains,
		PrivateKey: privateKey,
		Bundle:     true, // Returned byte[] will contain both issuer's and issued certificate
	}

	cert, err := client.Certificate.Obtain(request)
	if err != nil {
		return nil, err
	}

	return cert, nil
}
