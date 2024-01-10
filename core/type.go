package core

type ProvisionRequest struct {
	Attestation string `json:"attestation"`
	Domain      string `json:"domain"`
}
