package api

import "github.com/akakou/ra_webs/verifier/ent"

func checkTAValiditiy(servs []*ent.TAServer) bool {
	isValid := true
	for _, serv := range servs {
		isValid = isValid && len(serv.Edges.Violation) == 0
	}

	return isValid && len(servs) > 0 && servs[0].HasActivated
}
