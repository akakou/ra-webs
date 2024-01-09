module github.com/akakou/ra_webs/ttp/example

go 1.21.4

replace github.com/akakou/ra_webs/ta => ../

replace github.com/akakou/ra_webs/core => ../../core

require (
	github.com/akakou/ra_webs/core v0.0.0-00010101000000-000000000000 // indirect
	github.com/akakou/ra_webs/ta v0.0.0-00010101000000-000000000000 // indirect
)
