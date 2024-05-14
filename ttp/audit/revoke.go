package audit

import (
	"github.com/akakou/ra_webs/ttp/core"
	"github.com/akakou/ra_webs/ttp/ent"
	"github.com/akakou/ra_webs/ttp/ent/taserver"
)

func revoke(serv *ent.TAServer, db *core.DB) {
	db.Client.TAViolation.Create().
		SetServer(serv).
		SaveX(*db.Ctx)

	service, err := serv.QueryService().First(*db.Ctx)
	if err != nil {
		panic(err)
	}
	
	_, err = service.Update().SetIsActive(false).Save(*db.Ctx)
	if err != nil {
		panic(err)
	}
}

func revokeByDomain(domain string, last int, db *core.DB) {
	all, _ := db.Client.TAServer.
		Query().
		Where(taserver.DomainEQ(domain)).
		Where(taserver.IDGT(last - 1)).
		All(*db.Ctx)

	// todo: error handling

	for _, serv := range all {
		revoke(serv, db)
	}
}

func revokeByDomains(domains []string, db *core.DB) {
	for _, domain := range domains {
		last := lastValidID(domain, db)
		revokeByDomain(domain, last, db)
	}
}
