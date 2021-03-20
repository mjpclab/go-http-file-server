package tpl

import (
	"archive/zip"
	"errors"
	"html/template"
	"io"
	"net/http"
	"time"
)

type MemTheme struct {
	template *template.Template
	assets   assets
}

var initTime = time.Now()

func LoadMemTheme(themePath string) (theme MemTheme, err error) {
	var currentTheme = MemTheme{
		template: nil,
		assets:   make(assets, 2),
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
		if f.Name == templateFilename {
			currentTheme.template, err = ParsePageTpl(string(raw))
			if err != nil {
				return
			}
		} else {
			currentTheme.assets.set(f.Name, raw)
		}
	}

	if currentTheme.template != nil {
		theme = currentTheme
		return
	}

	err = errors.New("lacks of page template '" + templateFilename + "' in theme")
	return
}

func (theme MemTheme) RenderPage(w io.Writer, data interface{}) error {
	return theme.template.Execute(w, data)
}

func (theme MemTheme) RenderAsset(w http.ResponseWriter, r *http.Request, assetPath string) {
	asset, ok := theme.assets[assetPath]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", asset.contentType)
	http.ServeContent(w, r, assetPath, initTime, asset.readSeeker)
}
