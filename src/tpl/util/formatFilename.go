package util

import (
	"html/template"
	"strings"
)

var filenameReplacer = strings.NewReplacer(
	"\000", "<em>nul</em>",
	"\a", "<em>\\a</em>",
	"\b", "<em>\\b</em>",
	"\f", "<em>\\f</em>",
	"\n", "<em>\\n</em>",
	"\r", "<em>\\r</em>",
	"\t", "<em>\\t</em>",
	"\v", "<em>\\v</em>",

	"\000", "<em>\\000</em>",
	"\001", "<em>\\001</em>",
	"\002", "<em>\\002</em>",
	"\003", "<em>\\003</em>",
	"\004", "<em>\\004</em>",
	"\005", "<em>\\005</em>",
	"\006", "<em>\\006</em>",
	"\007", "<em>\\007</em>",
	"\010", "<em>\\010</em>",
	"\011", "<em>\\011</em>",
	"\012", "<em>\\012</em>",
	"\013", "<em>\\013</em>",
	"\014", "<em>\\014</em>",
	"\015", "<em>\\015</em>",
	"\016", "<em>\\016</em>",
	"\017", "<em>\\017</em>",
	"\020", "<em>\\020</em>",
	"\021", "<em>\\021</em>",
	"\022", "<em>\\022</em>",
	"\023", "<em>\\023</em>",
	"\024", "<em>\\024</em>",
	"\025", "<em>\\025</em>",
	"\026", "<em>\\026</em>",
	"\027", "<em>\\027</em>",
	"\030", "<em>\\030</em>",
	"\031", "<em>\\031</em>",
	"\032", "<em>\\032</em>",
	"\033", "<em>\\033</em>",
	"\034", "<em>\\034</em>",
	"\035", "<em>\\035</em>",
	"\036", "<em>\\036</em>",
	"\037", "<em>\\037</em>",
	"\177", "<em>\\177</em>",
)

func FormatFilename(filename string) template.HTML {
	escaped := template.HTMLEscapeString(filename)
	escaped = filenameReplacer.Replace(escaped)
	return template.HTML(escaped)
}
