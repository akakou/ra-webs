package ct

import (
	"github.com/akakou/ra_webs/ttp/core"
	"github.com/akakou/ra_webs/ttp/ent/taserver"
)

func logViolationByDomain(domain string, db *core.DB) error {
	serv, err := db.Client.TAServer.
		Query().
		Where(taserver.DomainEQ(domain)).
		Only(*db.Ctx)

	if err != nil {
		return nil
	}

	db.Client.TAViolation.Create().
		SetServer(serv).
		SaveX(*db.Ctx)

	service := serv.QueryService().FirstX(*db.Ctx)
	service.Update().SetIsActive(false).SaveX(*db.Ctx)

	return nil
}

func logViolationsByDomains(domains []string, db *core.DB) {
	for _, domain := range domains {
		logViolationByDomain(domain, db)
	}
}
