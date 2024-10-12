package direct

import (
	"context"
	"net/http"

	"errors"

	"github.com/akakou/ctstream/core"
	"github.com/akakou/ctstream/utils"
	"github.com/google/certificate-transparency-go/client"
	"github.com/google/certificate-transparency-go/jsonclient"
)

type CTClientParams struct {
	Start     core.LogID
	LogClient *client.LogClient
}

type CTClient struct {
	Url string
	*client.LogClient
	context.Context
	first        core.LogID
	opts         jsonclient.Options
	maxEntrySize int64
}

func NewCTClient(url string, maxEntrySize int64, ctx context.Context, ops jsonclient.Options) (*CTClient, error) {
	hc := http.Client{}
	_ctx := context.Context(ctx)

	c, err := client.New(url, &hc, ops)
	if err != nil {
		return nil, errors.New(ERROR_FAILED_TO_NEW)
	}

	return &CTClient{
		Url:          url,
		LogClient:    c,
		Context:      _ctx,
		maxEntrySize: maxEntrySize,
		opts:         ops,
	}, nil
}

func DefaultCTClient(url string, ctx context.Context) (*CTClient, error) {
	return NewCTClient(url, core.DefaultMaxEntries, ctx, jsonclient.Options{})
}

func (stream *CTClient) Init() error {
	sct, err := stream.GetSTH(stream.Context)
	if err != nil {
		return errors.New(ERROR_FAILED_TO_FETCH_STH)
	}

	stream.first = int64(sct.TreeSize)
	return nil
}

func (stream *CTClient) Next(callback core.Callback) {
	sct, err := stream.LogClient.GetSTH(stream.Context)
	if err != nil {
		err = errors.New(ERROR_FAILED_TO_FETCH_STH)
		callFailed(err, stream.LogClient, callback)
		return
	}

	treeSize := int64(sct.TreeSize)

	if treeSize == stream.first {
		return
	}

	var start int64
	var end int64

	for start = stream.first; end < treeSize; start += stream.maxEntrySize {
		end = start + stream.maxEntrySize

		if treeSize < end {
			end = treeSize
		}

		ThreadManager.Run(
			func(args Args) {
				start := args[0].(int64)
				end := args[1].(int64)

				entries, err := stream.LogClient.GetEntries(stream.Context, start, end)
				certs, err2 := extractCertFromEntries(entries)

				utils.Callbacks(certs, &CTClientParams{
					LogClient: stream.LogClient,
					Start:     start,
				}, callback, errors.Join(err, err2))
			}, Args{start, end},
		)
	}

	stream.first = treeSize
}

func (stream *CTClient) GetDomain() string {
	return stream.Url
}
