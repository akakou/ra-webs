package dns

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/akakou/ra_webs/core"
	"github.com/miekg/dns"
)

type Server struct {
	Zone    string
	storage RecordHolder
}

var DefaultZone = os.Getenv("ZONE")
var SelfIp = os.Getenv("SELF_IP")

func New() (*Server, error) {
	s := &Server{
		Zone:    DefaultZone,
		storage: NewInMemory(),
	}
	s.storage.Append(DefaultZone, SelfIp)
	return s, nil
}

func (s *Server) ToFqdn(hostname string) string {
	return fmt.Sprintf("%v.%v", hostname, s.Zone)
}

func (s *Server) AddHost(fqdn string, ip string) error {
	log.Printf("AddHost: Appending a host: %v => %v\n", fqdn, ip)
	return s.storage.Append(fqdn, ip)
}

func (s *Server) Lookup(fqdn string) (string, error) {
	queries := []string{
		fqdn,
		TrimTrailingPeriod(fqdn),
	}

	for _, q := range queries {
		ip, err := s.storage.Query(q)

		if err != nil {
			continue
		}

		return ip, nil
	}

	return "", ErrNotFound

}

func (s *Server) Serve(w dns.ResponseWriter, req *dns.Msg) {
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

	if errors.Is(err, ErrNotFound) || ipStr == "" {
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

	ar := &dns.A{
		Hdr: newRrHeader(query.Name),
		A:   ip,
	}

	caar := &dns.CAA{
		Hdr:   newRrHeader(query.Name),
		Flag:  0,
		Tag:   "issue",
		Value: core.CA_NAME,
	}

	m.Answer = append(m.Answer, ar)
	m.Answer = append(m.Answer, caar)
	w.WriteMsg(m)

	log.Printf("Lookup: served a query successfully: %v => %v\n", query.Name, ip)
}

func (s *Server) Start(addr string) error {
	dns.HandleFunc(s.Zone, s.Serve)

	server := &dns.Server{
		Addr: addr,
		Net:  "udp",
	}
	log.Printf("Start: Serving DNS queries at %v...\n", addr)
	return server.ListenAndServe()
}
