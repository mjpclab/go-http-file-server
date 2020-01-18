package util

import (
	"html/template"
	"strings"
)

var filenameReplacer = strings.NewReplacer(
	"\r", "<em>\\r</em>",
	"\n", "<em>\\n</em>",
	"\a", "<em>\\a</em>",
	"\v", "<em>\\v</em>",
)

func FormatFilename(filename string) template.HTML {
	escaped := template.HTMLEscapeString(filename)
	escaped = filenameReplacer.Replace(escaped)
	return template.HTML(escaped)
}
