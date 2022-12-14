package theme

import (
	"archive/zip"
	"errors"
	"html/template"
	"io"
	"net/http"
	"time"
)

type MemTheme struct {
	Template *template.Template
	Assets   Assets
}

var initTime = time.Now()

func LoadMemTheme(themePath string) (theme MemTheme, err error) {
	var currentTheme = MemTheme{
		Template: nil,
		Assets:   make(Assets, 3),
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
		return
	}

	err = errors.New("lacks of page template '" + templateFilename + "' in theme")
	return
}

func (theme MemTheme) RenderPage(w io.Writer, data interface{}) error {
	return theme.Template.Execute(w, data)
}

func (theme MemTheme) RenderAsset(w http.ResponseWriter, r *http.Request, assetPath string) {
	asset, ok := theme.Assets[assetPath]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", asset.ContentType)
	http.ServeContent(w, r, assetPath, initTime, asset.ReadSeeker)
}
