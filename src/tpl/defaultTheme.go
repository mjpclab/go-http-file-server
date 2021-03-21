package tpl

import (
	"./frontend"
	"strings"
)

var DefaultTheme MemTheme

func init() {
	defaultTpl, err := ParsePageTpl(frontend.DefaultTplStr)
	if err != nil {
		defaultTpl, _ = ParsePageTpl("Builtin Template Error")
	}
	DefaultTheme.template = defaultTpl

	defaultAssets := map[string]asset{
		"index.css": {"text/css", strings.NewReader(frontend.DefaultCss)},
		"index.js":  {"application/javascript", strings.NewReader(frontend.DefaultJs)},
	}
	DefaultTheme.assets = defaultAssets
}
