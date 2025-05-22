module github.com/akakou/ra-webs/monitor/serv/example

go 1.23.6

replace github.com/akakou/ra-webs/devkit/core => ../../../devkit/core/

replace github.com/akakou/ra-webs/devkit/service => ../../../devkit/service

replace github.com/akakou/ra-webs/monitor => ../../

replace github.com/akakou/ra-webs/monitor/serviceclient => ../../serviceclient

require (
	github.com/akakou/ra-webs/devkit/core v0.0.0-00010101000000-000000000000
	github.com/akakou/ra-webs/monitor v0.0.0-00010101000000-000000000000
	github.com/labstack/echo/v4 v4.13.3
)

require (
	github.com/SherClockHolmes/webpush-go v1.4.0 // indirect
	github.com/akakou/ctstream v0.0.39 // indirect
	github.com/akakou/extract-embed v0.0.1 // indirect
	github.com/akakou/ra-webs/devkit/service v0.0.0-00010101000000-000000000000 // indirect
	github.com/akakou/ra-webs/monitor/serviceclient v0.0.0-00010101000000-000000000000 // indirect
	github.com/akakou/sslmate-cert-search-api v0.0.15 // indirect
	github.com/apparentlymart/go-textseg/v15 v15.0.0 // indirect
	github.com/bmatcuk/doublestar v1.3.4 // indirect
	github.com/cockroachdb/errors v1.11.3 // indirect
	github.com/cockroachdb/logtags v0.0.0-20241215232642-bb51bb14a506 // indirect
	github.com/cockroachdb/redact v1.1.6 // indirect
	github.com/getsentry/sentry-go v0.31.1 // indirect
	github.com/go-jose/go-jose/v4 v4.0.5 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.2 // indirect
	github.com/google/certificate-transparency-go v1.2.1 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	github.com/zclconf/go-cty-yaml v1.1.0 // indirect
	golang.org/x/sync v0.12.0 // indirect
	golang.org/x/tools v0.31.0 // indirect
)

require (
	ariga.io/atlas v0.32.0 // indirect
	entgo.io/ent v0.14.4 // indirect
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/akakou/go-utils v0.0.7 // indirect
	github.com/akakou/golang-utils v0.0.1 // indirect
	github.com/edgelesssys/ego v1.7.0 // indirect
	github.com/go-openapi/inflect v0.21.2 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/uuid v1.6.0 // indirect; noi
	github.com/hashicorp/hcl/v2 v2.23.0 // indirect
	github.com/labstack/gommon v0.4.2 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-sqlite3 v1.14.27 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	github.com/zclconf/go-cty v1.16.2 // indirect
	golang.org/x/crypto v0.36.0 // indirect
	golang.org/x/mod v0.24.0 // indirect
	golang.org/x/net v0.38.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	golang.org/x/time v0.11.0 // indirect
)
