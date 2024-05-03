module github.com/akakou/ra-webs/domain-owner/example

go 1.21.9

replace github.com/akakou/ra-webs/domainowner => ../

require github.com/akakou/ra-webs/domainowner v0.0.0-00010101000000-000000000000

require (
	github.com/miekg/dns v1.1.59 // indirect
	golang.org/x/mod v0.16.0 // indirect
	golang.org/x/net v0.22.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	golang.org/x/tools v0.19.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
