package monitor

import (
	"github.com/akakou/ra_webs/verifier/core"
	"github.com/akakou/ra_webs/verifier/ent"
	"github.com/akakou/ra_webs/verifier/ent/taserver"
)

func revoke(serv *ent.TAServer, verifier *core.Verifier) {
	err := core.NotifierViolation(serv.Domain, verifier)
	if err != nil {
		panic(err)
	}

	verifier.DB.Client.TAViolation.Create().
		SetServer(serv).
		SaveX(*verifier.DB.Ctx)

	service, err := serv.QueryService().First(*verifier.DB.Ctx)
	if err != nil {
		panic(err)
	}

	_, err = service.Update().SetIsActive(false).Save(*verifier.DB.Ctx)
	if err != nil {
		panic(err)
	}
}

func revokeByDomain(domain string, last int, verifier *core.Verifier) {
	all, _ := verifier.DB.Client.TAServer.
		Query().
		Where(taserver.DomainEQ(domain)).
		Where(taserver.IDGT(last - 1)).
		All(*verifier.DB.Ctx)

	// todo: error handling

	for _, serv := range all {
		revoke(serv, verifier)
	}
}

func revokeByDomains(domains []string, verifier *core.Verifier) {
	for _, domain := range domains {
		last := lastValidID(domain, verifier.DB)
		revokeByDomain(domain, last, verifier)
	}
}
