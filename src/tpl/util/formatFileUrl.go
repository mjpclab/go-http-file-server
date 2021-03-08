package util

import (
	"strings"
)

var fileUrlReplacer = strings.NewReplacer(
	"%", "%25",
	"?", "%3f",
	"&", "%26",
	"#", "%23",
	"=", "%3d",
)

func FormatFileUrl(filename string) string {
	escaped := fileUrlReplacer.Replace(filename)
	return escaped
}
