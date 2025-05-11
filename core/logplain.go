package core

type LogPlain struct {
	Repository string `json:"repository"`
	CommitId   string `json:"commit_id"`
	Evidence   string `json:"evidence"`
	PublicKey  []byte `json:"public_key"`
}
