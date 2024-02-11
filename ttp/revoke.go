package ttp

import "github.com/akakou/ra_webs/ttp/ent/ta"

func revokeAllDomain(db *DB, domains []string) {
	for _, violatingDomain := range domains {
		taInfos, err := db.Client.TA.
			Query().
			Where(ta.DomainEQ(violatingDomain)).
			WithAuditLog().
			WithCode().
			All(*db.Ctx)

		if err != nil {
			continue
		}

		for _, taInfo := range taInfos {
			taInfo.Edges.AuditLog.IsValid = false
			taInfo.Edges.AuditLog.Update().Save(*db.Ctx)

		}
	}
}
