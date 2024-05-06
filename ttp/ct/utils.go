package ct

import (
	"crypto/x509"

	metact "github.com/akakou/meta-ct"
	"github.com/akakou/ra_webs/ttp/core"
	"github.com/akakou/ra_webs/ttp/ent"
	"github.com/akakou/ra_webs/ttp/ent/taserver"
)

func MetaCertsToCerts(cs []metact.MetaCert) ([]x509.Certificate, error) {
	certs := []x509.Certificate{}

	for _, c := range cs {
		cert, err := c.Certificate()

		if err != nil {
			return []x509.Certificate{}, err
		}

		certs = append(certs, *cert)
	}

	return certs, nil
}

func subscribeCT(domain string, ct *metact.MetaCT) error {
	return ct.Subscribe(domain)
}

var SubscribeCT = subscribeCT

func lastValidID(domain string, db *core.DB) int {
	lastValid, err := db.Client.TAServer.
		Query().
		Where(taserver.DomainEQ(domain)).
		Where(taserver.HasActivated(true)).
		Order(ent.Desc(taserver.FieldID)).
		First(*db.Ctx)

	var lastID = 0
	if err == nil {
		lastID = lastValid.ID
	}

	return lastID
}
