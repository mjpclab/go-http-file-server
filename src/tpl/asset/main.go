package asset

import (
	"io"
	"strings"
)

type content struct {
	ContentType string
	ReadSeeker  io.ReadSeeker
}

var assets = map[string]content{
	"main.css": {"text/css", strings.NewReader(mainCss)},
	"main.js":  {"application/javascript", strings.NewReader(mainJs)},
}

func Get(path string) (content, bool) {
	c, ok := assets[path]
	return c, ok
}
