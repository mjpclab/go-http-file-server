package tpl

import (
	"../serverErrHandler"
	"./util"
	_ "embed"
	"html/template"
	"path"
)

//go:embed frontend/page.html
var pageTplStr string

var defaultPageTpl *template.Template

func init() {
	pageTpl := template.New("page")
	pageTpl = addFuncMap(pageTpl)

	var err error
	defaultPageTpl, err = pageTpl.Parse(pageTplStr)
	if serverErrHandler.CheckError(err) {
		defaultPageTpl = template.Must(pageTpl.Parse("Builtin Template Error"))
	}
}

func LoadPageTpl(tplPath string) (*template.Template, error) {
	if len(tplPath) == 0 {
		return defaultPageTpl, nil
	}

	var err error
	pageTpl := template.New(path.Base(tplPath))
	pageTpl = addFuncMap(pageTpl)
	pageTpl, err = pageTpl.ParseFiles(tplPath)
	if err != nil {
		pageTpl = defaultPageTpl
	}

	return pageTpl, err
}

func addFuncMap(tpl *template.Template) *template.Template {
	return tpl.Funcs(template.FuncMap{
		"fmtFilename": util.FormatFilename,
		"fmtSize":     util.FormatSize,
		"fmtTime":     util.FormatTime,
	})
}
