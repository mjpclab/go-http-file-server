package util

import "testing"

func TestWildcardToStrRegexp(t *testing.T) {
	wildcard := "file.{name}*.txt"
	regexp := WildcardToStrRegexp(wildcard)
	if regexp != `^file\.\{name\}.*?\.txt$` {
		t.Error(regexp)
	}
}
