package direct

import (
	"context"
	"time"

	"github.com/akakou/ctstream/core"
)

var DefaultSleep = 1 * time.Second

func DefaultCTStream(url string, ctx context.Context) (*core.CTStream[*CTClient], error) {
	c, err := DefaultCTClient(url, ctx)
	if err != nil {
		return nil, err
	}

	return core.NewCTStream(
		c,
		DefaultSleep,
		ctx,
	)
}

func DefaultCTsStream(urls []string, ctx context.Context) (*core.PararellCTsStream[*core.CTStream[*CTClient]], error) {
	streams := []*core.CTStream[*CTClient]{}

	for _, url := range urls {
		stream, err := DefaultCTStream(url, ctx)
		if err != nil {
			return nil, err
		}

		streams = append(streams, stream)
	}

	return core.NewPararellCTsStream(streams, DefaultSleep)
}
