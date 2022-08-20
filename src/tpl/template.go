package tpl

import (
	"html/template"
	"mjpclab.dev/ghfs/src/tpl/util"
)

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
