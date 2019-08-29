package util

import "testing"

func TestWildcardToRegexp(t *testing.T) {
	wildcard := "file.{name}*.txt"
	regexp := WildcardToRegexp(wildcard)
	if regexp != `^file\.\{name\}.*?\.txt$` {
		t.Error(regexp)
	}
}
