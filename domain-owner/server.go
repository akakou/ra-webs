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

	switch req.Question[0].Qtype {
	case dns.TypeA:
		s.ServeLookUpA(w, m, req)
	case dns.TypeCAA:
		s.ServLookUpCAA(w, m, req)
	}
}

func (s *DNSServer) ServeLookUpA(w dns.ResponseWriter, m, req *dns.Msg) {
	query := req.Question[0]
	log.Printf("Lookup: query received: %#v\n", query)

	ipStr, err := s.Lookup(query.Name)

	if err != nil && strings.Contains(err.Error(), ERROR_NOT_FOUND) {
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
		Hdr: newAHeader(query.Name),
		A:   ip,
	}

	m.Answer = append(m.Answer, aRecord)
	w.WriteMsg(m)

	log.Printf("Lookup: served a query successfully: %v => %v\n\n", query.Name, ip)
}

func (s *DNSServer) ServLookUpCAA(w dns.ResponseWriter, m, req *dns.Msg) {
	query := req.Question[0]

	caaRecord1 := &dns.CAA{
		Hdr:   newCAAHeader(query.Name),
		Flag:  128, // it mean ca must not issue certificate if CA does not understand CAA
		Tag:   "issue",
		Value: CaName,
	}

	caaRecord2 := &dns.CAA{
		Hdr:   newCAAHeader(query.Name),
		Flag:  128, // it means that ca must not issue certificate if CA does not understand CAA
		Tag:   "issuewild",
		Value: ";", // it means no certificate having domain expressed by wild card not support
	}

	m.Answer = append(m.Answer, caaRecord1)
	m.Answer = append(m.Answer, caaRecord2)
	log.Printf("Lookup: served caa query successfully: %v => %v\n%v\n\n",
		query.Name,
		caaRecord1,
		caaRecord2,
	)

	w.WriteMsg(m)

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
