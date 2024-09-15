module github.com/akakou/ra_webs/verifier/example

go 1.21.4

replace github.com/akakou/ra_webs/ta => ../

replace github.com/akakou/ra_webs/core => ../../core

require github.com/akakou/ra_webs/ta v0.0.0-00010101000000-000000000000

require (
	github.com/akakou/ra_webs/core v0.0.0-00010101000000-000000000000 // indirect
	github.com/cenkalti/backoff/v4 v4.2.1 // indirect
	github.com/edgelesssys/ego v1.5.3 // indirect
	github.com/go-acme/lego/v4 v4.16.1 // indirect
	github.com/go-jose/go-jose/v4 v4.0.2 // indirect
	github.com/miekg/dns v1.1.58 // indirect
	golang.org/x/crypto v0.24.0 // indirect
	golang.org/x/mod v0.17.0 // indirect
	golang.org/x/net v0.25.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/sys v0.21.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	golang.org/x/tools v0.21.1-0.20240508182429-e35e4ccd0d2d // indirect
)
