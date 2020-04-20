package tpl

import (
	"../serverErrHandler"
	"./util"
	"html/template"
	"path"
)

const pageTplStr = `
<!DOCTYPE html>
<html lang="">
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
<meta http-equiv="X-UA-Compatible" content="IE=edge"/>
<meta name="viewport" content="initial-scale=1, user-scalable=no"/>
<meta name="format-detection" content="telephone=no"/>
<meta name="renderer" content="webkit"/>
<meta name="wap-font-scale" content="no"/>
<title>{{.Path}}</title>
<link rel="stylesheet" type="text/css" href="{{.RootRelPath}}?assert=main.css"/>
</head>
<body class="{{if .IsRoot}}root-dir{{else}}sub-dir{{end}}">
<ol class="path-list">
{{range .Paths}}
<li><a href="{{.Path}}">{{fmtFilename .Name}}</a></li>
{{end}}
</ol>
{{if .CanUpload}}
<div class="upload">
<form method="POST" action="{{.SubItemPrefix}}?upload" enctype="multipart/form-data">
<input type="file" name="file" multiple="multiple" class="file"/>
<input type="submit" value="Upload" class="submit"/>
</form>
</div>
{{end}}
{{if .CanMkdir}}
<div class="mkdir">
<form method="POST" action="{{.SubItemPrefix}}?mkdir">
<input type="text" name="name" autocomplete="off" class="name"/>
<input type="submit" value="mkdir" class="submit"/>
</form>
</div>
{{end}}
{{if .CanArchive}}
<div class="archive">
<a href="{{.SubItemPrefix}}?tar" download="{{.ItemName}}.tar">.tar</a>
<a href="{{.SubItemPrefix}}?tgz" download="{{.ItemName}}.tar.gz">.tar.gz</a>
<a href="{{.SubItemPrefix}}?zip" download="{{.ItemName}}.zip">.zip</a>
</div>
{{end}}
{{if .CanDelete}}
<script type="text/javascript">
function confirmDelete(el) {
var href = el.href;
var name = decodeURIComponent(href.substr(href.lastIndexOf('=') + 1));
return confirm('Delete?\n' + name);
}
</script>
{{end}}
<ul class="item-list{{if .HasDeletable}} has-deletable{{end}}">
<li class="dir parent">
<a href="{{if .IsRoot}}./{{else}}../{{end}}" class="link">
<span class="name">../</span>
<span class="size"></span>
<span class="time"></span>
</a>
</li>
{{range .SubItemsHtml}}
<li class="{{.Type}}">
<a href="{{.Url}}" class="link">
<span class="name">{{.DisplayName}}</span>
<span class="size">{{.DisplaySize}}</span>
<span class="time">{{.DisplayTime}}</span>
</a>
{{if .DeleteUrl}}<a href="{{.DeleteUrl}}" class="delete" onclick="return confirmDelete(this)">x</a>{{end}}
</li>
{{end}}
</ul>
{{if eq .Status 403}}
<div class="error">403 resource is forbidden</div>
{{else if eq .Status 404}}
<div class="error">404 resource not found</div>
{{else if eq .Status 500}}
<div class="error">500 potential issue occurred</div>
{{end}}
<script type="text/javascript" src="{{.RootRelPath}}?assert=main.js" defer="defer" async="async"></script>
</body>
</html>
`

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
