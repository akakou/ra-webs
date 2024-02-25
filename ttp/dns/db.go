package dns

import "github.com/akakou/ra_webs/ttp/db"

func (s *Server) LoadDatabase(db *db.DB) error {
	tas, err := db.Client.TA.Query().WithServer().All(*db.Ctx)
	if err != nil {
		return err
	}

	for _, ta := range tas {
		err := s.AddHost(ta.Edges.Server.Domain, ta.Edges.Server.IP)
		if err != nil {
			return err
		}
	}

	return nil
}
