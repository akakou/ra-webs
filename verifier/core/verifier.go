package core

import "github.com/labstack/echo/v4"

type Verifier struct {
	DB         *DB
	Monitor    Monitor
	Notifier   Notifier
	AdminToken string
}

func NewVerifier(db *DB, monitor Monitor, notifier Notifier, adminToken string) (*Verifier, error) {
	return &Verifier{
		DB:         db,
		Monitor:    monitor,
		AdminToken: adminToken,
		Notifier:   notifier,
	}, nil
}

func (verifier *Verifier) Setup(e *echo.Group) error {
	err := verifier.Monitor.Setup(verifier)
	if err != nil {
		return err
	}

	err = verifier.Notifier.Setup(e, verifier)
	if err != nil {
		return err
	}

	return nil
}
