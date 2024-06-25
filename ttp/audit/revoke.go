package audit

import (
	"github.com/akakou/ra_webs/ttp/core"
	"github.com/akakou/ra_webs/ttp/ent"
	"github.com/akakou/ra_webs/ttp/ent/taserver"
)

func revoke(serv *ent.TAServer, ttp *core.TTP) {
	err := core.NotifyViolation(serv.Domain, ttp)
	if err != nil {
		panic(err)
	}

	ttp.DB.Client.TAViolation.Create().
		SetServer(serv).
		SaveX(*ttp.DB.Ctx)

	service, err := serv.QueryService().First(*ttp.DB.Ctx)
	if err != nil {
		panic(err)
	}

	_, err = service.Update().SetIsActive(false).Save(*ttp.DB.Ctx)
	if err != nil {
		panic(err)
	}
}

func revokeByDomain(domain string, last int, ttp *core.TTP) {
	all, _ := ttp.DB.Client.TAServer.
		Query().
		Where(taserver.DomainEQ(domain)).
		Where(taserver.IDGT(last - 1)).
		All(*ttp.DB.Ctx)

	// todo: error handling

	for _, serv := range all {
		revoke(serv, ttp)
	}
}

func revokeByDomains(domains []string, ttp *core.TTP) {
	for _, domain := range domains {
		last := lastValidID(domain, ttp.DB)
		revokeByDomain(domain, last, ttp)
	}
}
