package domainowner

import "github.com/miekg/dns"

const DefaultTtl = 120

func newRrHeader(hostname string) dns.RR_Header {
	return dns.RR_Header{
		Name:   hostname,
		Rrtype: dns.TypeA,
		Class:  dns.ClassINET,
		Ttl:    DefaultTtl,
	}
}
