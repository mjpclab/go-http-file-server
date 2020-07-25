package util

import (
	"html/template"
	"time"
)

func FormatTime(t time.Time) template.HTML {
	return template.HTML(t.Format("2006-01-02<span> 15:04</span>"))
}
