package theme

import (
	"bytes"
	"io"
	"mjpclab.dev/ghfs/src/util"
)

type Asset struct {
	ContentType string
	ReadSeeker  io.ReadSeeker
}

type Assets map[string]Asset

func (assets Assets) Set(path string, content []byte) error {
	rd := bytes.NewReader(content)
	ctype, err := util.GetContentType(path, rd)
	if err != nil {
		return err
	}

	asset := Asset{
		ContentType: ctype,
		ReadSeeker:  rd,
	}
	assets[path] = asset
	return nil
}
