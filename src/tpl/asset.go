package tpl

import (
	"../util"
	"bytes"
	_ "embed"
	"io"
)

type asset struct {
	ContentType string
	ReadSeeker  io.ReadSeeker
}

type assets map[string]asset

func (assets assets) Set(path string, content []byte) error {
	rd := bytes.NewReader(content)
	ctype, err := util.GetContentType(path, rd)
	if err != nil {
		return err
	}

	asset := asset{
		ContentType: ctype,
		ReadSeeker:  rd,
	}
	assets[path] = asset
	return nil
}

func (assets assets) Get(path string) (asset, bool) {
	c, ok := assets[path]
	return c, ok
}

//go:embed frontend/index.css
var css []byte

//go:embed frontend/index.js
var js []byte

var defaultAssets = map[string]asset{
	"index.css": {"text/css", bytes.NewReader(css)},
	"index.js":  {"application/javascript", bytes.NewReader(js)},
}
