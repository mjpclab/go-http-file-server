package util

import (
	"runtime"
	"testing"
)

func TestWildcardToStrRegexp(t *testing.T) {
	wildcard := "file.{name}*.txt"
	regexp := WildcardToStrRegexp(wildcard)
	expected := `^file\.\{name\}.*?\.txt$`
	if runtime.GOOS == "windows" {
		expected = "(?i)" + expected
	}
	if regexp != expected {
		t.Error(regexp + " <= " + expected)
	}
}
