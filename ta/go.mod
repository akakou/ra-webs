module github.com/akakou/ra_webs/ta

go 1.21.4

replace github.com/akakou/ra_webs/core => ../core

require (
	github.com/akakou/ra_webs/core v0.0.0-00010101000000-000000000000
	github.com/edgelesssys/ego v1.4.1
	github.com/go-playground/assert/v2 v2.2.0
)

require (
	github.com/cenkalti/backoff/v4 v4.2.1 // indirect
	github.com/go-acme/lego/v4 v4.14.2 // indirect
	github.com/go-jose/go-jose/v3 v3.0.0 // indirect
	github.com/miekg/dns v1.1.55 // indirect
	golang.org/x/crypto v0.12.0 // indirect
	golang.org/x/mod v0.11.0 // indirect
	golang.org/x/net v0.11.0 // indirect
	golang.org/x/sys v0.11.0 // indirect
	golang.org/x/text v0.12.0 // indirect
	golang.org/x/tools v0.10.0 // indirect
	gopkg.in/square/go-jose.v2 v2.6.0 // indirect
)
