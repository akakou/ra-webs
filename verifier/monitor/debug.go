package monitor

import (
	"context"

	"github.com/akakou/ctstream/direct"
)

func DefaultDirectMonitor() (*Monitor, error) {
	ctx := context.Background()
	stream, err := direct.DefaultCTsStream(DefaultCTLogs, ctx)
	if err != nil {
		return nil, err
	}

	return NewMonitor(stream)
}
