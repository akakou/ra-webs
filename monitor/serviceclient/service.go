package serviceclient

type EvidenceEntry struct {
	Repository string `json:"repository"`
	CommitID   string `json:"commit_id"`
	Evidence   string `json:"evidence"`
}

type ServiceClient interface {
	Fetch(publicKey []byte) (*EvidenceEntry, error)
}
