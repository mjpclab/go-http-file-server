package defaultTheme

import (
	"bytes"
	_ "embed"
	"mjpclab.dev/ghfs/src/tpl/theme"
)

//go:embed frontend/index.html
var defaultTplStr string

//go:embed frontend/index.css
var defaultCss []byte

//go:embed frontend/index.js
var defaultJs []byte

//go:embed frontend/favicon.ico
var defaultFavicon []byte

var DefaultTheme theme.MemTheme

func init() {
	var err error

	DefaultTheme.Template, err = theme.ParsePageTpl(defaultTplStr)
	if err != nil {
		DefaultTheme.Template, _ = theme.ParsePageTpl("Builtin Template Error")
	}

	DefaultTheme.Assets = theme.Assets{
		{"index.css", "text/css; charset=utf-8", bytes.NewReader(defaultCss)},
		{"index.js", "application/javascript; charset=utf-8", bytes.NewReader(defaultJs)},
		{"favicon.ico", "image/x-icon", bytes.NewReader(defaultFavicon)},
	}
}
