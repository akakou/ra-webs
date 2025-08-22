package service

import "errors"

var (
	ErrNoToken        = errors.New("no token")
	ErrNoDomain       = errors.New("no domain")
	ErrNoATPrivateKey = errors.New("not private key")
	ErrFailedToOpenDB = errors.New("failed to open DB")
	ErrCreateSchema   = errors.New("failed to create schema")
	ErrCreateAudit    = errors.New("failed to create audit")
)
