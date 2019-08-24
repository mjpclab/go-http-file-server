package tpl

import (
	"../fmtSize"
	"../serverError"
	"path"
	"text/template"
)

const pageTplStr = `
<!DOCTYPE html>
<html>
<head>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
    <meta http-equiv="X-UA-Compatible" content="IE=edge"/>
    <meta name="viewport" content="initial-scale=1"/>
    <meta name="format-detection" content="telephone=no"/>
    <meta name="renderer" content="webkit"/>
    <meta name="wap-font-scale" content="no"/>
    <base href="{{.Scheme}}//{{.Host}}/{{if .Path}}{{.Path}}/{{end}}"/>

    <title>{{.Path}}</title>

    <style type="text/css">
        html, body {
            margin: 0;
            padding: 0;
            background: #fff;
        }

        html {
            font-family: "roboto_condensedbold", "Helvetica Neue", Helvetica, Arial, sans-serif;
        }

        body {
            color: #333;
            font-size: 0.625em;
            font-family: Consolas, "DejaVu Sans Mono", Monospaced;
        }

        a {
            display: block;
            padding: 0.25em 0.5em;
            color: inherit;
            text-decoration: none;
        }

        a:hover {
            color: #000;
            background: #f5f5f5;
        }

        .path-list {
            font-size: 1.4em;
            overflow: hidden;
            border-bottom: 1px #999 solid;
        }

        .path-list a {
            position: relative;
            float: left;
            padding-right: 1.2em;
            text-align: center;
            white-space: nowrap;
            min-width: 1em;
        }

        .path-list a:after {
            content: '';
            position: absolute;
            top: 0.6em;
            right: 0.4em;
            width: 0.4em;
            height: 0.4em;
            border: 1px solid;
            border-color: #ccc #ccc transparent transparent;
            transform: rotate(45deg);
        }

        .path-list a:last-child {
            padding-right: 0.5em;
        }

        .path-list a:last-child:after {
            display: none;
        }

        .item-list {
            padding: 1em;
        }

        .item-list a {
            display: flex;
            flex-flow: row nowrap;
            align-items: center;
            border-bottom: 1px #f5f5f5 solid;
        }

        .item-list span {
            margin-left: 1em;
            flex-shrink: 0;
        }

        .item-list .name {
            flex: 1 1 0;
            margin-left: 0;
            font-size: 1.4em;
            word-break: break-all;
        }

        .item-list .size {
            white-space: nowrap;
            text-align: right;
            color: #666;
        }

        .item-list .time {
            width: 10em;
            color: #999;
            text-align: right;
            white-space: nowrap;
            overflow: hidden;
        }

        .error {
            margin: 1em;
            padding: 1em;
            background: #ffc;
        }
    </style>
</head>
<body>

<div class="path-list">
    <a href="/">/</a>
    {{range .Paths}}
        {{with .}}
            <a href="{{.Path}}">{{.Name}}</a>
        {{end}}
    {{end}}
</div>

<div class="item-list">
    <a href="../">
        <span class="name">../</span>
        <span class="size"></span>
        <span class="time"></span>
    </a>
    {{range .SubItems}}
        {{$isDir := .IsDir}}
        <a href="{{.Name}}" class="item {{if $isDir}}item-dir{{else}}item-file{{end}}">
            <span class="name">{{.Name}}{{if $isDir}}/{{end}}</span>
            <span class="size">{{if not $isDir}}{{fmtSize .Size}}{{end}}</span>
            <span class="time">{{printf "%04d-%02d-%02d %02d:%02d" .ModTime.Year .ModTime.Month .ModTime.Day .ModTime.Hour .ModTime.Minute}}</span>
        </a>
    {{end}}
</div>

{{if .Error}}
    <div class="error">{{.Error}}</div>
{{end}}

</body>
</html>
`

var defaultPage *template.Template

func init() {
	tplObj := template.New("page")
	tplObj = addFuncMap(tplObj)

	var err error
	defaultPage, err = tplObj.Parse(pageTplStr)
	if serverError.CheckError(err) {
		defaultPage = template.Must(tplObj.Parse("Builtin Template Error"))
	}
}

func LoadPage(tplPath string) *template.Template {
	var tplObj *template.Template
	var err error

	if len(tplPath) > 0 {
		tplObj = template.New(path.Base(tplPath))
		tplObj = addFuncMap(tplObj)
		tplObj, err = tplObj.ParseFiles(tplPath)
		serverError.CheckError(err)
	}
	if err != nil || len(tplPath) == 0 {
		tplObj = defaultPage
	}

	return tplObj
}

func addFuncMap(tpl *template.Template) *template.Template {
	return tpl.Funcs(template.FuncMap{
		"fmtSize": fmtSize.FmtSize,
	})
}
