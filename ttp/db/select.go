package db

import (
	"github.com/akakou/ra_webs/ttp/ent"
	"github.com/akakou/ra_webs/ttp/ent/ta"
	"github.com/akakou/ra_webs/ttp/ent/taserver"
)

func (db *DB) SelectTAByDomain(domain string) (*ent.TA, error) {
	serv, err := db.Client.TAServer.Query().Where(taserver.DomainEQ(domain)).First(*db.Ctx)
	if err != nil {
		return nil, err
	}

	return serv.QueryTa().WithCode().WithCtAudit().WithServer().First(*db.Ctx)
}

func (db *DB) SelectTAByService(service *ent.Service) (*ent.TA, error) {
	ta, err := service.QueryTaserver().QueryTa().WithCode().WithCtAudit().WithServer().First(*db.Ctx)
	if err != nil {
		return nil, err
	}

	return ta, nil
}

func (db *DB) SelectTA(id int) (*ent.TA, error) {
	return db.Client.TA.Query().Where(ta.ID(id)).WithCode().WithCtAudit().WithServer().First(*db.Ctx)
}
