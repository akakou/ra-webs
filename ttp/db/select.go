package db

import (
	"github.com/akakou/ra_webs/ttp/ent"
	"github.com/akakou/ra_webs/ttp/ent/taserver"
)

func SelectTAByDomain(domain string, db *DB) (*ent.TA, error) {
	serv, err := db.Client.TAServer.Query().Where(taserver.DomainEQ(domain)).First(*db.Ctx)
	if err != nil {
		return nil, err
	}

	return serv.QueryTa().WithCode().WithCtAudit().WithServer().First(*db.Ctx)
}
