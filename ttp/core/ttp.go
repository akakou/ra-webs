package core

import "github.com/labstack/echo/v4"

type TTP struct {
	DB         *DB
	CT         CT
	AdminToken string
}

func NewTTP(db *DB, ct CT, adminToken string) (*TTP, error) {
	return &TTP{
		DB:         db,
		CT:         ct,
		AdminToken: adminToken,
	}, nil
}

func (ttp *TTP) Setup(e *echo.Echo) error {
	return ttp.CT.Setup(e, ttp)
}
