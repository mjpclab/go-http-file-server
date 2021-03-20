package tpl

import (
	"io"
	"net/http"
)

const templateFilename = "index.html"

type Theme interface {
	RenderPage(w io.Writer, data interface{}) error
	RenderAsset(w http.ResponseWriter, r *http.Request, assetPath string)
}
