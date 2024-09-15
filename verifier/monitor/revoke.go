package monitor

import (
	"github.com/akakou/ra_webs/verifier/core"
	"github.com/akakou/ra_webs/verifier/ent"
)

func revoke(serv *ent.TAServer, verifier *core.Verifier) {
	err := core.NotifierViolation(serv.Domain, verifier)
	if err != nil {
		panic(err)
	}

	verifier.DB.Client.TAViolation.Create().
		SetServer(serv).
		SaveX(*verifier.DB.Ctx)

	service, err := serv.QueryService().Only(*verifier.DB.Ctx)
	if err != nil {
		panic(err)
	}

	_, err = service.Update().SetIsActive(false).Save(*verifier.DB.Ctx)
	if err != nil {
		panic(err)
	}
}
