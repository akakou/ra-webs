package api

import (
	"errors"

	verifiercore "github.com/akakou/ra_webs/verifier/core"
	"github.com/akakou/ra_webs/verifier/ent"
)

func checkTAValidity(serv *ent.TAServer, verifier *verifiercore.Verifier) (bool, error) {
	code, err := serv.QueryCode().WithService().Only(*verifier.DB.Ctx)
	if err != nil {
		return false, errors.New("code is not found")
	}

	service, err := serv.QueryService().Only(*verifier.DB.Ctx)
	if err != nil {
		return false, errors.New("service is not found")
	}
	return serv.HasActivated && service.IsActive && code.IsActive && code.Edges.Service.IsActive, nil
}

func checkViolationLogs(servs []*ent.TAServer) (bool, error) {
	isValid := true

	if len(servs) == 0 {
		return false, errors.New("server is not found")
	}

	for _, serv := range servs {
		isValid = isValid && len(serv.Edges.Violation) == 0
	}
	return isValid, nil
}
