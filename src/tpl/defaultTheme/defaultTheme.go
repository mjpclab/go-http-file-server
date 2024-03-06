package defaultTheme

import (
	"bytes"
	"mjpclab.dev/ghfs/src/tpl/defaultTheme/frontend"
	"mjpclab.dev/ghfs/src/tpl/theme"
	"strings"
)

var DefaultTheme theme.MemTheme

func init() {
	var err error

	DefaultTheme.Template, err = theme.ParsePageTpl(frontend.DefaultTplStr)
	if err != nil {
		DefaultTheme.Template, _ = theme.ParsePageTpl("Builtin Template Error")
	}

	DefaultTheme.Assets = theme.Assets{
		{"index.css", "text/css; charset=utf-8", strings.NewReader(frontend.DefaultCss)},
		{"index.js", "application/javascript; charset=utf-8", strings.NewReader(frontend.DefaultJs)},
		{"favicon.ico", "image/x-icon", bytes.NewReader(frontend.DefaultFavicon)},
	}
}
