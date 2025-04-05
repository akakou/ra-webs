package logclient

import (
	"github.com/akakou/ra-webs/log/api/interfacestruct"
)

func (logclient *LogClient) Fetch() (*interfacestruct.TA, error) {
	result := interfacestruct.TA{
		Evidence:   "",
		Signature:  []byte(""),
		Repository: "",
		CommitID:   "",
	}

	return &result, nil
}
