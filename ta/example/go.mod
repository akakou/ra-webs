module github.com/akakou/ra_webs/ttp/example

go 1.21.4

replace github.com/akakou/ra_webs/ta => ../

replace github.com/akakou/ra_webs/core => ../../core

require (
	github.com/akakou/go-utils v0.0.3
	github.com/akakou/ra_webs/core v0.0.0-00010101000000-000000000000
	github.com/akakou/ra_webs/ta v0.0.0-00010101000000-000000000000
	github.com/labstack/echo/v4 v4.12.0
)

require (
	github.com/cenkalti/backoff/v4 v4.2.1 // indirect
	github.com/edgelesssys/ego v1.5.0 // indirect
	github.com/go-acme/lego/v4 v4.16.1 // indirect
	github.com/go-jose/go-jose/v4 v4.0.1 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/labstack/gommon v0.4.2 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/miekg/dns v1.1.58 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	golang.org/x/crypto v0.22.0 // indirect
	golang.org/x/mod v0.14.0 // indirect
	golang.org/x/net v0.24.0 // indirect
	golang.org/x/sys v0.19.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/time v0.5.0 // indirect
	golang.org/x/tools v0.17.0 // indirect
	gopkg.in/square/go-jose.v2 v2.6.0 // indirect
)
