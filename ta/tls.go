package ta

import (
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

const CERT_DIER_CACHE = "./tmp/ra-webs.cache"

func (ap *TA) TLSConfig() (autocert.Manager, error) {
	acmeClient := acme.Client{
		DirectoryURL: AcmeURL,
	}

	return autocert.Manager{
		Client:   &acmeClient,
		Cache:    autocert.DirCache(CERT_DIER_CACHE),
		Prompt:   autocert.AcceptTOS,
		ForceRSA: true,
	}, nil
}
