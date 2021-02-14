package util

import (
	"testing"
)

func TestFormatFileUrl(t *testing.T) {
	if FormatFileUrl("abc") != "abc" {
		t.Error()
	}

	if FormatFileUrl("a?b") != "a%3fb" {
		t.Error()
	}

	if FormatFileUrl("a&b") != "a%26b" {
		t.Error()
	}

	if FormatFileUrl("a#b") != "a%23b" {
		t.Error()
	}

	if FormatFileUrl("a=b") != "a%3db" {
		t.Error()
	}
}
