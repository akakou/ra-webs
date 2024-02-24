package ttp

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	metact "github.com/akakou/meta-ct"
	"github.com/akakou/ra_webs/ttp/ent"
	"github.com/akakou/ra_webs/ttp/ent/service"
	"github.com/akakou/ra_webs/ttp/ent/taserver"
	"github.com/labstack/echo/v4"
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

func revokeByDomain(db *DB, domains []string) {
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

func authenticateService(db *DB, c echo.Context) (*ent.Service, error) {
	authorization := c.Request().Header["Authorization"][0]
	token := authorization[len("Bearer "):]

	service, err := db.Client.Service.Query().Where(service.TokenEQ(token)).First(*db.Ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate service: %w", err)
	}

	return service, nil
}

func authenticateAdmin(auditor *Auditor, c echo.Context) error {
	authorization := c.Request().Header["Authorization"][0]
	token := authorization[len("Bearer "):]

	if token != auditor.adminToken {
		return c.String(http.StatusUnauthorized, "token is invalid")
	}

	return nil
}
