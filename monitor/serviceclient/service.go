package serviceclient

type EvidenceEntry struct {
	Repository string
	CommitID   string
	Evidence   string
	PublicKey  []byte
}

type ServiceClient interface {
	Fetch(publicKey []byte) (*EvidenceEntry, error)
}
