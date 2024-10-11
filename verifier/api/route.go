package api

import (
	verifiercore "github.com/akakou/ra_webs/verifier/core"
	"github.com/labstack/echo/v4"
)

func Route(e *echo.Group, verifier *verifiercore.Verifier) {
	RegisterApi.Set(e, verifier)
	GetServerApi.Set(e, verifier)

	PostServiceByAdmin.Set(e, verifier)

	GetServerFromDomainApi.Set(e, verifier)
	PostNotifierApi.Set(e, verifier)

	PostReloadDBApi.Set(e, verifier)
}
