package tpl

import (
	"../serverErrHandler"
	"../util"
	"html/template"
	"path"
)

const pageTplStr = `
{{$subItemPrefix := .SubItemPrefix}}
<!DOCTYPE html>
<html lang="">
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
<meta http-equiv="X-UA-Compatible" content="IE=edge"/>
<meta name="viewport" content="initial-scale=1"/>
<meta name="format-detection" content="telephone=no"/>
<meta name="renderer" content="webkit"/>
<meta name="wap-font-scale" content="no"/>
<title>{{.Path}}</title>
<link rel="stylesheet" type="text/css" href="{{.RootRelPath}}?assert=main.css"/>
</head>
<body>
<ol class="path-list">
{{range .Paths}}
<li><a href="{{.Path}}">{{html .Name}}</a></li>
{{end}}
</ol>
{{if .CanUpload}}
<div class="upload">
<form method="POST" enctype="multipart/form-data">
<input type="file" name="files" class="files" multiple="multiple" accept="*/*"/>
<input type="submit" value="Upload"/>
</form>
</div>
{{end}}
{{if .CanArchive}}
<div class="archive">
<a href="{{$subItemPrefix}}?tar" download="{{.ItemName}}.tar">.tar</a>
<a href="{{$subItemPrefix}}?tgz" download="{{.ItemName}}.tar.gz">.tar.gz</a>
<a href="{{$subItemPrefix}}?zip" download="{{.ItemName}}.zip">.zip</a>
</div>
{{end}}
<ul class="item-list">
<li>
<a href="{{if .IsRoot}}./{{else}}../{{end}}">
<span class="name">../</span>
<span class="size"></span>
<span class="time"></span>
</a>
</li>
{{range .SubItems}}
{{$isDir := .IsDir}}
<li>
<a href="{{$subItemPrefix}}{{.Name}}{{if $isDir}}/{{end}}" class="item {{if $isDir}}item-dir{{else}}item-file{{end}}">
<span class="name">{{html .Name}}{{if $isDir}}/{{end}}</span>
<span class="size">{{if not $isDir}}{{fmtSize .Size}}{{end}}</span>
<span class="time">{{fmtTime .ModTime}}</span>
</a>
</li>
{{end}}
</ul>
{{range .Errors}}
<div class="error">{{.}}</div>
{{end}}
<script type="text/javascript" src="{{.RootRelPath}}?assert=main.js"></script>
</body>
</html>
`

var defaultPage *template.Template

func init() {
	tplObj := template.New("page")
	tplObj = addFuncMap(tplObj)

	var err error
	defaultPage, err = tplObj.Parse(pageTplStr)
	if serverErrHandler.CheckError(err) {
		defaultPage = template.Must(tplObj.Parse("Builtin Template Error"))
	}
}

func LoadPage(tplPath string) (*template.Template, error) {
	var tplObj *template.Template
	var err error

	if len(tplPath) > 0 {
		tplObj = template.New(path.Base(tplPath))
		tplObj = addFuncMap(tplObj)
		tplObj, err = tplObj.ParseFiles(tplPath)
	}
	if err != nil || len(tplPath) == 0 {
		tplObj = defaultPage
	}

	return tplObj, err
}

func addFuncMap(tpl *template.Template) *template.Template {
	return tpl.Funcs(template.FuncMap{
		"fmtSize": util.FormatSize,
		"fmtTime": util.FormatTimeMinute,
	})
}
