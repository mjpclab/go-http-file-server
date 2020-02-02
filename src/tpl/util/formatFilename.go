package util

import (
	"html/template"
	"strings"
)

var filenameReplacer = strings.NewReplacer(
	"\a", "<em>\\a</em>",
	"\b", "<em>\\b</em>",
	"\f", "<em>\\f</em>",
	"\n", "<em>\\n</em>",
	"\r", "<em>\\r</em>",
	"\t", "<em>\\t</em>",
	"\v", "<em>\\v</em>",
)

func FormatFilename(filename string) template.HTML {
	escaped := template.HTMLEscapeString(filename)
	escaped = filenameReplacer.Replace(escaped)
	return template.HTML(escaped)
}
