package theme

import (
	"bytes"
	"io"
	"mjpclab.dev/ghfs/src/util"
)

type Asset struct {
	Path        string
	ContentType string
	ReadSeeker  io.ReadSeeker
}

type Assets []Asset

func (assets Assets) Append(path string, content []byte) (Assets, error) {
	rd := bytes.NewReader(content)
	ctype, err := util.GetContentType(path, rd)
	if err != nil {
		return assets, err
	}

	assets = append(assets, Asset{
		Path:        path,
		ContentType: ctype,
		ReadSeeker:  rd,
	})

	return assets, nil
}
