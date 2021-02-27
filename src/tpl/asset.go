package tpl

import (
	"bytes"
	_ "embed"
	"io"
)

//go:embed frontend/index.css
var css []byte

//go:embed frontend/index.js
var js []byte

type content struct {
	ContentType string
	ReadSeeker  io.ReadSeeker
}

var assets = map[string]content{
	"index.css": {"text/css", bytes.NewReader(css)},
	"index.js":  {"application/javascript", bytes.NewReader(js)},
}

func GetAsset(path string) (content, bool) {
	c, ok := assets[path]
	return c, ok
}
