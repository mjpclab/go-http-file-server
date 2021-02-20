package tpl

import (
	"./frontend"
	"io"
	"strings"
)

type content struct {
	ContentType string
	ReadSeeker  io.ReadSeeker
}

var assets = map[string]content{
	"main.css": {"text/css", strings.NewReader(frontend.MainCss)},
	"main.js":  {"application/javascript", strings.NewReader(frontend.MainJs)},
}

func GetAsset(path string) (content, bool) {
	c, ok := assets[path]
	return c, ok
}
