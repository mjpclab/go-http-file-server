package tpl

import (
	"archive/zip"
	"html/template"
	"io"
)

type Theme struct {
	Template *template.Template
	Assets   assets
}

var DefaultTheme Theme

// wait for `defaultTpl` initialized
func init() {
	DefaultTheme = Theme{defaultTpl, defaultAssets}
}

func LoadTheme(themePath string) (theme Theme, err error) {
	theme = DefaultTheme

	if len(themePath) == 0 {
		return
	}

	var currentTheme = Theme{
		Template: nil,
		Assets:   make(assets, 2),
	}
	// assume to be a zip file
	var zipRd *zip.ReadCloser
	zipRd, err = zip.OpenReader(themePath)
	if err != nil {
		return
	}
	defer zipRd.Close()

	for _, f := range zipRd.File {
		var rd io.ReadCloser
		rd, err = f.Open()
		if err != nil {
			continue
		}
		var raw []byte
		raw, err = io.ReadAll(rd)
		rd.Close()
		if err != nil {
			return
		}
		if f.Name == "index.html" {
			currentTheme.Template, err = ParsePageTpl(string(raw))
			if err != nil {
				return
			}
		} else {
			currentTheme.Assets.Set(f.Name, raw)
		}
	}

	if currentTheme.Template != nil {
		theme = currentTheme
	}

	return
}
