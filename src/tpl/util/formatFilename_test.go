package util

import (
	"html/template"
	"testing"
)

func TestFormatFilename(t *testing.T) {
	raw := "a\rb\nc"
	replaced := FormatFilename(raw)
	if replaced != template.HTML("a<em>\\r</em>b<em>\\n</em>c") {
		t.Error(replaced)
	}
}
