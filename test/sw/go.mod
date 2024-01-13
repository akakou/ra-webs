module github.com/akakou/ra_webs/test

go 1.21.4

replace github.com/akakou/ra_webs/core => ../../core

replace github.com/akakou/ra_webs/ta => ../../ta

replace github.com/akakou/ra_webs/ttp => ../../ttp

require github.com/akakou/ra_webs/ta v0.0.0-00010101000000-000000000000

require (
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/cenkalti/backoff/v4 v4.2.1 // indirect
	github.com/edgelesssys/ego v1.4.1 // indirect
	github.com/go-acme/lego v2.7.2+incompatible // indirect
	github.com/go-acme/lego/v4 v4.14.2 // indirect
	github.com/go-jose/go-jose/v3 v3.0.0 // indirect
	github.com/miekg/dns v1.1.55 // indirect
	golang.org/x/tools v0.10.0 // indirect
	gopkg.in/square/go-jose.v2 v2.6.0 // indirect
)

require (
	github.com/akakou/ra_webs/core v0.0.0-00010101000000-000000000000 // indirect
	github.com/labstack/echo/v4 v4.11.4
	github.com/labstack/gommon v0.4.2 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	golang.org/x/crypto v0.17.0 // indirect
	golang.org/x/mod v0.11.0 // indirect
	golang.org/x/net v0.19.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
	golang.org/x/text v0.14.0 // indirect
)
