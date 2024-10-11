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

	service := serv.QueryService().OnlyX(*verifier.DB.Ctx)

	service.Update().SetIsActive(false).SaveX(*verifier.DB.Ctx)
}
