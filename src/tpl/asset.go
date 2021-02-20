package tpl

import (
	"bytes"
	_ "embed"
	"io"
)

//go:embed frontend/main.css
var mainCss []byte

//go:embed frontend/main.js
var mainJs []byte

type content struct {
	ContentType string
	ReadSeeker  io.ReadSeeker
}

var assets = map[string]content{
	"main.css": {"text/css", bytes.NewReader(mainCss)},
	"main.js":  {"application/javascript", bytes.NewReader(mainJs)},
}

func GetAsset(path string) (content, bool) {
	c, ok := assets[path]
	return c, ok
}
