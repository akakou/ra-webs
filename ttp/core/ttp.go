package core

import "github.com/labstack/echo/v4"

type TTP struct {
	DB         *DB
	Audit      Audit
	Notify     Notify
	AdminToken string
}

func NewTTP(db *DB, audit Audit, notify Notify, adminToken string) (*TTP, error) {
	return &TTP{
		DB:         db,
		Audit:      audit,
		AdminToken: adminToken,
		Notify:     notify,
	}, nil
}

func (ttp *TTP) Setup(e *echo.Echo) error {
	err := ttp.Audit.Setup(ttp)
	if err != nil {
		return err
	}

	err = ttp.Notify.Setup(e, ttp)
	if err != nil {
		return err
	}

	return nil
}
