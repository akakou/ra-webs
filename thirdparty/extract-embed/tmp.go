package extractembed

import (
	"embed"
	"os"
)

type TmpExtractedFs struct {
	Path string
}

func ExtractFs(name string, embedFs *embed.FS) (*TmpExtractedFs, error) {
	dname, err := os.MkdirTemp("", name)
	if err != nil {
		return nil, err
	}

	err = Extract(dname, embedFs)
	if err != nil {
		return nil, err
	}

	return &TmpExtractedFs{
		Path: dname,
	}, nil
}

func (t *TmpExtractedFs) Close() error {
	return os.RemoveAll(t.Path)
}
