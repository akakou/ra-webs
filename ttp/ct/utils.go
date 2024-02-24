package ct

import (
	"errors"
	"strings"

	metact "github.com/akakou/meta-ct"
	"github.com/akakou/ra_webs/ttp/core"
	"github.com/akakou/ra_webs/ttp/ent/taserver"
)

func findCertExtensions(extensions []metact.KeyValue, label string) (string, error) {
	for _, ext := range extensions {
		if ext.Key == label {
			return ext.Value, nil
		}
	}

	return "", errors.New("extension not found")
}

func extractDomainLast(domain string) string {
	domain = strings.Replace(domain, "*", "", -1)
	splited := strings.Split(domain, ".")

	var indexInt int
	if len(splited) >= 2 {
		indexInt = 2
	} else {
		indexInt = 1
	}

	last := splited[len(splited)-indexInt:]
	lastDomain := strings.Join(last, ".")

	return lastDomain
}

func revokeByDomain(db *core.DB, domains []string) {
	for _, violatingDomain := range domains {
		taServer, err := db.Client.TAServer.
			Query().
			WithTa().
			Where(taserver.DomainEQ(violatingDomain)).
			First(*db.Ctx)

		if err != nil {
			continue
		}

		ta := taServer.Edges.Ta
		ta.IsValid = false
		ta.Update().Save(*db.Ctx)
	}
}
