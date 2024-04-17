package ct

import metact "github.com/akakou/meta-ct"

func debugSubscribeCT(domain string, ct *metact.MetaCT) error {
	return nil
}

func EnableDebug() {
	SubscribeCT = debugSubscribeCT
}
