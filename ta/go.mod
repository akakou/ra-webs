module github.com/akakou/ra_webs/ta

go 1.21.4

replace github.com/akakou/ra_webs/core => ../core

require (
	github.com/akakou/ra_webs/core v0.0.0-00010101000000-000000000000
	github.com/edgelesssys/ego v1.4.1
	github.com/go-playground/assert/v2 v2.2.0
)

require (
	golang.org/x/crypto v0.12.0 // indirect
	gopkg.in/square/go-jose.v2 v2.6.0 // indirect
)
