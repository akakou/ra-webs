package ct

import (
	"errors"

	metact "github.com/akakou/meta-ct"
	"github.com/akakou/ra_webs/ttp/db"
	"github.com/akakou/ra_webs/ttp/ent/taserver"
)

func findCertExtensions(extensions []metact.KeyValue, label string) (string, error) {
	for _, ext := range extensions {
		if ext.Key == label {
			return ext.Value, nil
		}
	}

	return "", errors.New("extension not found")
}

func revokeTAbyDomain(domain string, db *db.DB) error {
	taServer, err := db.Client.TAServer.
		Query().
		WithTa().
		Where(taserver.DomainEQ(domain)).
		Only(*db.Ctx)

	if err != nil {
		return nil
	}

	for _, ta := range taServer.Edges.Ta {
		ta.IsValid = false
		_, err := ta.Update().Save(*db.Ctx)

		if err != nil {
			return err
		}
	}

	return nil
}

func revokeTAByDomains(db *db.DB, domains []string) {
	for _, domain := range domains {
		revokeTAbyDomain(domain, db)
	}
}
