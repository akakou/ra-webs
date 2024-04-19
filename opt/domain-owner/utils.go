package domainowner

import (
	"fmt"
	"strings"
)

func trimTrailingPeriod(fqdn string) string {
	s := fqdn
	s = strings.TrimSuffix(s, ".")

	return s
}

func toFqdn(hostname, zone string) string {
	return fmt.Sprintf("%v.%v", hostname, zone)
}
