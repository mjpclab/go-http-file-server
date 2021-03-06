package tpl

import (
	"../serverErrHandler"
	"./util"
	_ "embed"
	"html/template"
)

//go:embed frontend/index.html
var defaultTplStr string

var defaultTpl *template.Template

func init() {
	tpl := template.New("page")
	tpl = addFuncMap(tpl)

	var err error
	defaultTpl, err = tpl.Parse(defaultTplStr)
	if serverErrHandler.CheckError(err) {
		defaultTpl = template.Must(tpl.Parse("Builtin Template Error"))
	}
}

func ParsePageTpl(tplText string) (tpl *template.Template, err error) {
	tpl = template.New("page")
	tpl = addFuncMap(tpl)
	tpl, err = tpl.Parse(tplText)

	return
}

func addFuncMap(tpl *template.Template) *template.Template {
	return tpl.Funcs(template.FuncMap{
		"fmtFilename": util.FormatFilename,
		"fmtSize":     util.FormatSize,
		"fmtTime":     util.FormatTime,
	})
}
