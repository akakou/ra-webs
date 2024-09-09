package api

import "github.com/akakou/ra_webs/verifier/ent"

func checkTAValiditiy(servs []*ent.TAServer) bool {
	isValid := true
	for _, serv := range servs {
		isValid = isValid && len(serv.Edges.Violation) == 0
	}

	return isValid && servs[0].HasActivated && len(servs) > 0
}
