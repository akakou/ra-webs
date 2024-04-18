package domainowner

import "strings"

func TrimTrailingPeriod(fqdn string) string {
	s := fqdn
	s = strings.TrimSuffix(s, ".")

	return s
}
