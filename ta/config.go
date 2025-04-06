package ta

import (
	"fmt"
	"os"

	"github.com/akakou/ra-webs/core"
)

type TAConfig struct {
	Token       string
	Repository  string
	CommitID    string
	TADomain    string
	ATLogDomain string
	Email       string
}

func DefaultConfig() (*TAConfig, error) {
	token := os.Getenv("RA_WEBS_SERVICE_TOKEN")
	repository := os.Getenv("RA_WEBS_TA_REPOSITORY")
	commitId := os.Getenv("RA_WEBS_TA_COMMIT_ID")
	domain := os.Getenv("RA_WEBS_TA_DOMAIN")
	email := os.Getenv("RA_WEBS_SERVICE_EMAIL")
	atLogBaseEnv := os.Getenv("RA_WEBS_AT_LOG_BASE")

	if token == "" {
		return nil, fmt.Errorf("%v", ERROR_TOKEN_NOT_SET)
	}

	if repository == "" {
		return nil, fmt.Errorf("%v", ERROR_REPOSITORY_NOT_SET)
	}

	if email == "" {
		return nil, fmt.Errorf("%v", ERROR_EMAIL_NOT_SET)
	}

	if atLogBaseEnv == "" {
		atLogBaseEnv = "http://localhost" + core.MonitorPort
		fmt.Printf("RA_WEBS_AT_LOG_BASE is not set: so use %v\n", atLogBaseEnv)
	}

	if domain == "" {
		domain = "http://localhost" + core.TAPort
		fmt.Printf("RA_WEBS_TA_DOMAIN is not set: so use %v\n", domain)
	}

	return &TAConfig{
		Token:       token,
		Repository:  repository,
		CommitID:    commitId,
		TADomain:    domain,
		ATLogDomain: atLogBaseEnv,
	}, nil
}
