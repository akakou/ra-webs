package interfacestruct

type PostRequest struct {
	Repository string `json:"repository"`
	CommitId   string `json:"commit_id"`
	Evidence   string `json:"evidence"`
}
