package assert

import (
	"io"
	"strings"
)

type content struct {
	ContentType string
	ReadSeeker  io.ReadSeeker
}

var asserts = map[string]content{
	"main.css": {"text/css", strings.NewReader(mainCss)},
	"main.js":  {"application/javascript", strings.NewReader(mainJs)},
}

func Get(path string) (content, bool) {
	c, ok := asserts[path]
	return c, ok
}
