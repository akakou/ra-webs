package ta

import (
	"crypto"
	"crypto/rsa"
	"crypto/tls"
	"fmt"
	"log"

	"github.com/go-acme/lego/challenge/tlsalpn01"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
)

type ACMEConfig struct {
	Email      string
	Domain     string
	PrivateKey *rsa.PrivateKey
}

type acmeContext struct {
	user   *acmeUser
	client *lego.Client
	config *ACMEConfig
}

func initACMEContext(config *ACMEConfig) (*acmeContext, error) {
	acmeUser := &acmeUser{
		Email: config.Email,
	}

	legoConfig := lego.NewConfig(acmeUser)
	legoClient, err := lego.NewClient(legoConfig)

	if err != nil {
		return nil, fmt.Errorf("InitACMEContext: creating ACME Client: %v", err)
	}

	return &acmeContext{
		config: config,
		user:   acmeUser,
		client: legoClient,
	}, nil
}

func (ctx *acmeContext) issueCert() (*tls.Certificate, error) {
	err := ctx.setupCRProvider()
	if err != nil {
		return nil, fmt.Errorf("issueCert: setting up provider: %v", err)
	}

	err = ctx.register(ctx.client)
	if err != nil {
		return nil, fmt.Errorf("issueCert: registering acme user: %v", err)
	}

	resource, err := ctx.obtainCertificate()
	if err != nil {
		return nil, fmt.Errorf("issueCert: obtaining certificate: %v", err)
	}

	log.Printf("issueCert: obtained certificate for domain: %v\n", resource.Domain)

	cert := tls.Certificate{
		Certificate: [][]byte{resource.Certificate},
	}

	return &cert, nil
}

func (ctx *acmeContext) setupCRProvider() error {
	server := tlsalpn01.NewProviderServer("", "443")
	provider := ctx.client.Challenge.SetTLSALPN01Provider(server)
	return provider
}

func (ctx *acmeContext) register(client *lego.Client) error {
	options := registration.RegisterOptions{
		TermsOfServiceAgreed: true,
	}

	reg, err := client.Registration.Register(options)
	if err != nil {
		return err
	}

	ctx.user.Registration = reg
	return nil
}

func (ctx *acmeContext) obtainCertificate() (*certificate.Resource, error) {
	request := certificate.ObtainRequest{
		Domains:    []string{ctx.config.Domain},
		PrivateKey: ctx.config.PrivateKey,
		Bundle:     true, // Returned byte[] will contain both issuer's and issued certificate
	}

	cert, err := ctx.client.Certificate.Obtain(request)
	if err != nil {
		return nil, err
	}

	return cert, nil
}

type acmeUser struct {
	Email        string
	key          crypto.PrivateKey
	Registration *registration.Resource
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
