module github.com/akakou/ra-webs/ta

go 1.21.4

replace github.com/akakou/ra-webs/core => ../core

require (
	github.com/akakou/ra-webs/core v0.0.0-00010101000000-000000000000
	github.com/go-acme/lego/v4 v4.16.1
	golang.org/x/crypto v0.22.0
)

require (
	github.com/cenkalti/backoff/v4 v4.2.1 // indirect
	github.com/edgelesssys/ego v1.5.0 // indirect
	github.com/go-jose/go-jose/v4 v4.0.1 // indirect
	github.com/miekg/dns v1.1.58 // indirect
	golang.org/x/mod v0.14.0 // indirect
	golang.org/x/net v0.24.0 // indirect
	golang.org/x/sys v0.19.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/tools v0.17.0 // indirect
	gopkg.in/square/go-jose.v2 v2.6.0 // indirect
)
