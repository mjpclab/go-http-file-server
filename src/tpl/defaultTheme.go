package tpl

import (
	"bytes"
	_ "embed"
)

//go:embed frontend/index.html
var defaultTplStr string

//go:embed frontend/favicon.ico
var defaultFavicon []byte

//go:embed frontend/index.css
var defaultCss []byte

//go:embed frontend/index.js
var defaultJs []byte

var DefaultTheme MemTheme

func init() {
	defaultTpl, err := ParsePageTpl(defaultTplStr)
	if err != nil {
		defaultTpl, _ = ParsePageTpl("Builtin Template Error")
	}
	DefaultTheme.template = defaultTpl

	defaultAssets := map[string]asset{
		"favicon.ico": {"image/x-icon", bytes.NewReader(defaultFavicon)},
		"index.css":   {"text/css", bytes.NewReader(defaultCss)},
		"index.js":    {"application/javascript", bytes.NewReader(defaultJs)},
	}
	DefaultTheme.assets = defaultAssets
}
