package audit

import (
	"github.com/akakou/ra_webs/ttp/core"
	"github.com/akakou/ra_webs/ttp/ent"
	"github.com/akakou/ra_webs/ttp/ent/taserver"
)

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
