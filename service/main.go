package service

import (
	"os"

	goutils "github.com/akakou/go-utils"
)

type Service struct {
	Token   string
	TTPBase string
}

func NewService(token, ttp string) *Service {
	return &Service{Token: token, TTPBase: ttp}
}

func DefaultService() *Service {
	ttpBase := goutils.GetEnv("RA_WEBS_TTP_BASE", "localhost:8081")
	token := os.Getenv("RA_WEBS_SERVICE_TOKEN")

	if token == "" {
		panic("RA_WEBS_SERVICE_TOKEN is not set")
	}

	return NewService(token, ttpBase)
}
