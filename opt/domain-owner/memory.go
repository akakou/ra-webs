package domainowner

import (
	"log"
	"strings"
)

type DNSRecords map[string]string

func NewDNSRecords() *DNSRecords {
	return &DNSRecords{}
}

func (s DNSRecords) AppendFQDN(fqdn string, ip string) {
	log.Printf("Append: Appending a host: %v => %v\n", fqdn, ip)

	lower := strings.ToLower(fqdn)
	s[lower] = ip
}

func (s DNSRecords) AppendDomain(domain string, ip string, zone string) {
	fqdn := toFqdn(domain, zone)
	s.AppendFQDN(fqdn, ip)
}

func (s DNSRecords) Query(fqdn string) (string, error) {
	lower := strings.ToLower(fqdn)
	ip, ok := s[lower]
	if !ok {
		return "", ErrNotFound
	}

	return ip, nil
}

func (s DNSRecords) FromDomains(dict map[string]string, zone string) {
	for domain, ip := range dict {
		s.AppendDomain(domain, ip, zone)
	}
}
