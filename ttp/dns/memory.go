package dns

import "strings"

type RecordHolder interface {
	Append(fqdn string, ip string) error
	Query(fqdn string) (string, error)
}

type InMemory map[string]string

func NewInMemory() InMemory {
	return InMemory{}
}

func (s InMemory) Append(fqdn string, ip string) error {
	lower := strings.ToLower(fqdn)
	s[lower] = ip
	return nil
}

func (s InMemory) Query(fqdn string) (string, error) {
	lower := strings.ToLower(fqdn)
	ip, ok := s[lower]
	if !ok {
		return "", ErrNotFound
	}

	return ip, nil
}
