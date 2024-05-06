package ta

const (
	ERROR_DEFAULT_CONFIG     = "failed to read config"
	ERROR_GENERATE_RSA_KEY   = "failed to generate rsa key"
	ERROR_ATTEST_PUBLIC_KEY  = "failed to attest public key"
	ERROR_REQUEST_FAILED     = "request failed"
	ERROR_READ_BODY          = "read body failed"
	ERROR_STATUS_NOT_OK      = "status not ok"
	ERROR_TTP_BASE_PARSE     = "ttp base parse failed"
	ERROR_REPOSITORY_NOT_SET = "failed to read RA_WEBS_TA_REPOSITORY"
	ERROR_TOKEN_NOT_SET      = "failed to read RA_WEBS_SERVICE_TOKEN"
	ERROR_DOMAIN_NOT_SET     = "failed to read RA_WEBS_TA_DOMAIN"
	ERROR_EMAIL_NOT_SET      = "failed to read RA_WEBS_TA_EMAIL"
	ERROR_TTP_BASE_NOT_SET   = "failed to read RA_WEBS_TTP_BASE"
)

// token := os.Getenv("RA_WEBS_SERVICE_TOKEN")
// repository := os.Getenv("RA_WEBS_TA_REPOSITORY")
// domain := os.Getenv("RA_WEBS_TA_DOMAIN")
// email := os.Getenv("RA_WEBS_TA_EMAIL")
// ttpBase := os.Getenv("RA_WEBS_TTP_BASE")
