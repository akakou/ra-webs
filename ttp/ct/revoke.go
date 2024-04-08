package ct

import (
	"github.com/akakou/ra_webs/ttp/db"
	"github.com/akakou/ra_webs/ttp/ent"
	"github.com/akakou/ra_webs/ttp/ent/taserver"
)

func revokeTA(t *ent.TA, db *db.DB) {
	t.IsValid = false
	t.Update().SaveX(*db.Ctx)
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
		revokeTA(ta, db)
	}

	return nil
}

func revokeTAByDomains(db *db.DB, domains []string) {
	for _, domain := range domains {
		revokeTAbyDomain(domain, db)
	}
}
