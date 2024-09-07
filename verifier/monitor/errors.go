package monitor

const (
	ERROR_SELECT_TA       = "failed to select ta info from db"
	ERROR_SELECT_LAST_LOG = "failed to select last log from db"
	ERROR_SELECT_TA_CODE  = "failed to select ta code from db"
	ERROR_SELECT_SERVER   = "failed to select server from db"

	ERROR_LAST_TA_INVALID = "last TA is invalid"
	ERROR_DOMAIN_INVALID  = "domain is invalid"

	ERROR_DOMAIN_INVALID_BY_WILDCARD                  = "wildcard domain is not allowed"
	ERROR_DOMAIN_INVALID_BY_NUM_DOMAIN                = "number of domain must be 1"
	ERROR_DOMAIN_INVALID_NOT_MATCH_COMMONNAME_AND_SAT = "CN and SAT must be same"

	ERROR_CERTIFICATE_NOT_FOUND = "certificate not match"

	ERROR_EXTENSION_NOT_FOUND  = "extension not found"
	ERROR_PUBLIC_KEY_NOT_RSA   = "public key is not RSA"
	ERROR_PUBLIC_KEY_NOT_MATCH = "public key is not match"
)
