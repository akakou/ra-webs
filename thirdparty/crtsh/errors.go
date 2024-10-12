package crtsh

import "github.com/pkg/errors"

var (
	ErrorParseRssUrl      = errors.Errorf("failed to parse Parse RSS URL")
	ErrorParseIdUrl       = errors.Errorf("failed to parse Parse ID URL")
	ErrorParseCertificate = errors.Errorf("failed to parse Certificate")
	ErrorParsePem         = errors.Errorf("failed to parse PEM")
	ErrorParseHtml        = errors.Errorf("failed to parse HTML")
	ErrorParseInt         = errors.Errorf("failed to parse Int")
	ErrorFetchRss         = errors.Errorf("failed to fetch RSS feed")
)
