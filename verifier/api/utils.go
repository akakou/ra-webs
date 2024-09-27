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

	if !serv.IsActive {
		return false, errors.New("server has been not activated")
	}

	if !code.IsActive {
		return false, errors.New("code is not active")
	}

	if !service.IsActive {
		return false, errors.New("service of server is not active")
	}

	if !code.Edges.Service.IsActive {
		return false, errors.New("service of code is not active")
	}

	return true, nil
}

func checkViolationLogs(servs []*ent.TAServer) (bool, error) {
	if len(servs) == 0 {
		return false, errors.New("server is not found")
	}

	for _, serv := range servs {
		if len(serv.Edges.Violation) != 0 {
			return false, errors.New("server has violation logs")
		}
	}

	return true, nil
}
