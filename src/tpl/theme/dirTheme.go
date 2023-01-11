package theme

import (
	"io"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
)

type DirTheme string

func (dir DirTheme) RenderPage(w io.Writer, data interface{}) error {
	filename := string(dir) + "/" + templateFilename
	tplStr, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	tpl, err := ParsePageTpl(string(tplStr))
	if err != nil {
		return err
	}

	tpl.Execute(w, data)
	return nil
}

func (dir DirTheme) RenderAsset(w http.ResponseWriter, r *http.Request, assetPath string) {
	header := w.Header()
	header.Set("Cache-Control", "public, max-age=0")
	filename := string(dir) + "/" + strings.Replace(path.Clean(assetPath), "../", "", -1)
	http.ServeFile(w, r, filename)
}
