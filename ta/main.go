package ta

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

const CERT_DIER_CACHE = "/var/www/.cache"
const ATTEST_ENDPOINT = "/rawebs/attest"

type REGISTER_TYPE int

const (
	TTP REGISTER_TYPE = iota + 1
	ACME
)

func (ap *TA) TLSConfig(e *echo.Echo) error {
	taId, err := ap.Register()
	if err != nil {
		return err
	}

	switch ap.Config.Type {
	case TTP:
		err = ap.IssueTTPCert(taId, e)
	case ACME:
		err = ap.IssueAcmeCert(e)
	default:
		err = fmt.Errorf("unknown register type: %d", ap.Config.Type)
	}

	return err
}
