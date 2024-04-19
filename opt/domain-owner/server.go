package domainowner

import (
	"errors"
	"log"
	"net"
	"os"
	"strings"

	"github.com/miekg/dns"
)

type DNSServer struct {
	Zone    string
	Records *DNSRecords
}

var DefaultZone = os.Getenv("ZONE")
var SelfIp = os.Getenv("SELF_IP")
var CaName = os.Getenv("CA_NAME")

func NewDNSServer() (*DNSServer, error) {
	s := &DNSServer{
		Zone:    DefaultZone,
		Records: NewDNSRecords(),
	}

	s.Records.AppendFQDN(DefaultZone, SelfIp)
	return s, nil
}

func (s *DNSServer) Lookup(fqdn string) (string, error) {
	queries := []string{
		fqdn,
		trimTrailingPeriod(fqdn),
	}

	for _, q := range queries {
		ip, err := s.Records.Query(q)

		if err != nil {
			continue
		}

		return ip, nil
	}

	return "", errors.New(ERROR_NOT_FOUND)

}

func (s *DNSServer) Serve(w dns.ResponseWriter, req *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(req)
	m.Authoritative = true

	if len(req.Question) != 1 {
		log.Printf("Lookup: error: invalid length of query: %v\n", len(req.Question))
		return
	}

	query := req.Question[0]
	log.Printf("Lookup: query received: %#v\n", query)

	ipStr, err := s.Lookup(query.Name)

	if strings.Contains(err.Error(), ERROR_NOT_FOUND) || ipStr == "" {
		// host not found
		log.Printf("Lookup: host not found: %v\n", query.Name)
		m.SetRcode(req, dns.RcodeNameError)
		w.WriteMsg(m)
		return
	}

	if err != nil {
		log.Printf("Lookup: error: failed to query for host '%v': %v\n", query.Name, err)
		return
	}

	ip := net.ParseIP(ipStr)

	aRecord := &dns.A{
		Hdr: newRrHeader(query.Name),
		A:   ip,
	}

	caaRecord1 := &dns.CAA{
		Hdr:   newRrHeader(query.Name),
		Flag:  128, // it mean ca must not issue certificate if CA does not understand CAA
		Tag:   "issue",
		Value: CaName,
	}

	caaRecord2 := &dns.CAA{
		Hdr:   newRrHeader(query.Name),
		Flag:  128, // it means that ca must not issue certificate if CA does not understand CAA
		Tag:   "issuewild",
		Value: ";", // it means no certificate having domain expressed by wild card not support
	}

	m.Answer = append(m.Answer, aRecord)
	m.Answer = append(m.Answer, caaRecord1)
	m.Answer = append(m.Answer, caaRecord2)
	w.WriteMsg(m)

	log.Printf("Lookup: served a query successfully: %v => %v\n", query.Name, ip)
}

func (s *DNSServer) Start(addr string) error {
	dns.HandleFunc(s.Zone, s.Serve)

	server := &dns.Server{
		Addr: addr,
		Net:  "udp",
	}
	log.Printf("Start: Serving DNS queries at %v...\n", addr)
	return server.ListenAndServe()
}

func (s *DNSServer) AddHost(fqdn string, ip string) {
	s.Records.AppendFQDN(fqdn, ip)
}
