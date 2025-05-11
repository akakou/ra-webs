package monitor

import "errors"

var (
	errPublicKeyNotRSA   = errors.New("public key is not RSA")
	errPublicKeyNotMatch = errors.New("public key does not match the server public key")

	errInitDB       = errors.New("failed to initialize DB")
	errOpenDB       = errors.New("failed to open DB")
	errCreateSchema = errors.New("failed to create schema in DB")
	errCreateAudit  = errors.New("failed to create monitor")

	errDomainEnvironmentVariableIsEmpty = errors.New("Environment variable is empty")
)
