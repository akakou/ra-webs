package core

import (
	"time"

	ctx509 "github.com/google/certificate-transparency-go/x509"
)

type LogID = int64

var DefaultMaxEntries int64 = 256

var DefaultEpochSleep = time.Minute * 20
var DefaultPullingSleep = time.Second * 10

type Callback func(*ctx509.Certificate, int, any, error)

type CtClient interface {
	Init() error
	Next(callback Callback)
	GetDomain() string
}
