package ta

import (
	"fmt"
	"os"

	"github.com/akakou/ra-webs/devkit/core"
)

type TAConfig struct {
	Token       string
	Repository  string
	CommitID    string
	TADomain    string
	ServiceBase string
	Email       string
}

func DefaultConfig() (*TAConfig, error) {
	token := os.Getenv("RA_WEBS_SERVICE_TOKEN")
	repository := os.Getenv("RA_WEBS_TA_REPOSITORY")
	commitId := os.Getenv("RA_WEBS_TA_COMMIT_ID")
	domain := os.Getenv("RA_WEBS_TA_DOMAIN")
	email := os.Getenv("RA_WEBS_SERVICE_EMAIL")
	serviceBase := os.Getenv("RA_WEBS_SERVICE_BASE")

	if token == "" {
		return nil, fmt.Errorf("%v", ERROR_TOKEN_NOT_SET)
	}

	if repository == "" {
		return nil, fmt.Errorf("%v", ERROR_REPOSITORY_NOT_SET)
	}

	if email == "" {
		return nil, fmt.Errorf("%v", ERROR_EMAIL_NOT_SET)
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
		ServiceBase: serviceBase,
	}, nil
}
