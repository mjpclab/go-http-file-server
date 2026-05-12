package defaultTheme

import (
	_ "embed"
	"strings"

	"mjpclab.dev/ghfs/src/tpl/theme"
)

//go:embed frontend/index.html
var defaultTplStr string

//go:embed frontend/index.css
var defaultCss string

//go:embed frontend/index.js
var defaultJs string

//go:embed frontend/favicon.ico
var defaultFavicon string

var DefaultTheme theme.MemTheme

func init() {
	var err error

	DefaultTheme.Template, err = theme.ParsePageTpl(defaultTplStr)
	if err != nil {
		DefaultTheme.Template, _ = theme.ParsePageTpl("Builtin Template Error")
	}

	DefaultTheme.Assets = theme.Assets{
		{"index.css", "text/css; charset=utf-8", strings.NewReader(defaultCss)},
		{"index.js", "application/javascript; charset=utf-8", strings.NewReader(defaultJs)},
		{"favicon.ico", "image/x-icon", strings.NewReader(defaultFavicon)},
	}
}
