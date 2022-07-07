package serverHandler

import "testing"

func TestStripUrlPrefix(t *testing.T) {
	var result string

	result = stripUrlPrefix("/", "/", "/")
	if result != "/" {
		t.Error(result)
	}

	result = stripUrlPrefix("/foo", "/foo", "/")
	if result != "/foo" {
		t.Error(result)
	}

	result = stripUrlPrefix("/foo/bar", "/foo/bar", "/")
	if result != "/foo/bar" {
		t.Error(result)
	}

	result = stripUrlPrefix("/foo/bar/", "/foo/bar", "/")
	if result != "/foo/bar/" {
		t.Error(result)
	}

	result = stripUrlPrefix("/foo", "/foo", "/foo")
	if result != "/" {
		t.Error(result)
	}

	result = stripUrlPrefix("/foo/", "/foo", "/foo")
	if result != "/" {
		t.Error(result)
	}

	result = stripUrlPrefix("/foo/bar", "/foo/bar", "/foo")
	if result != "/bar" {
		t.Error(result)
	}

	result = stripUrlPrefix("/foo/bar/", "/foo/bar", "/foo")
	if result != "/bar/" {
		t.Error(result)
	}
}
