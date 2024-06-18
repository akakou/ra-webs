package core

import "github.com/labstack/echo/v4"

type TTP struct {
	DB         *DB
	Audit      Audit
	AdminToken string
}

func NewTTP(db *DB, audit Audit, adminToken string) (*TTP, error) {
	return &TTP{
		DB:         db,
		Audit:      audit,
		AdminToken: adminToken,
	}, nil
}

func (ttp *TTP) Setup(e *echo.Echo) error {
	return ttp.Audit.Setup(ttp)
}
