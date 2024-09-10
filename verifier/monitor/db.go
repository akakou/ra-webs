package monitor

import (
	"context"

	ctcore "github.com/akakou/ctstream/core"
	"github.com/akakou/ctstream/thirdparty/sslmate"
	"github.com/akakou/ra_webs/verifier/core"
	"github.com/akakou/ra_webs/verifier/ent"
	"github.com/akakou/ra_webs/verifier/ent/taserver"
)

func LoadSSLMateMonitor(db *core.DB, ctx context.Context) (*SSLMateMonitor, error) {
	servers, err := db.Client.TAServer.Query().Select(taserver.FieldDomain).All(*db.Ctx)

	if err != nil {
		return nil, err
	}

	return sslMateMonitorFromServers(servers, ctx)
}

func sslMateMonitorFromServers(servers []*ent.TAServer, ctx context.Context) (*SSLMateMonitor, error) {
	var streams []*ctcore.CTStream[*sslmate.SSLMateCTClient]

	for _, server := range servers {
		client, err := sslmate.DefaultCTClient(server.Domain)
		if err != nil {
			return nil, err
		}

		client.First = server.LastCtlog

		stream, err := sslmate.NewCTStream(client, -1, ctx)

		if err != nil {
			return nil, err
		}

		streams = append(streams, stream)
	}

	stream, err := sslmate.NewCTsStream(streams, sslmate.DefaultEpochSleep)

	return &SSLMateMonitor{
		ctstream: stream,
		ctx:      ctx,
	}, err
}

func UpdateLastLog(last string, domain string, db *core.DB) error {
	serv, err := db.Client.TAServer.
		Query().
		Where(taserver.DomainEQ(domain)).
		Order(ent.Desc(taserver.FieldID)).
		First(*db.Ctx)

	if err != nil {
		return err
	}

	serv.LastCtlog = last

	_, err = serv.Update().Save(*db.Ctx)

	return err
}
