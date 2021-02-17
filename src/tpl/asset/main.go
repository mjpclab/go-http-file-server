package asset

import (
	"bytes"
	"io"
)

type content struct {
	ContentType string
	ReadSeeker  io.ReadSeeker
}

var assets = map[string]content{
	"main.css": {"text/css", bytes.NewReader(mainCss)},
	"main.js":  {"application/javascript", bytes.NewReader(mainJs)},
}

func Get(path string) (content, bool) {
	c, ok := assets[path]
	return c, ok
}
