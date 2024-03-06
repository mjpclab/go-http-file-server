package util

import "testing"

func TestHasPrefixDirAccurate(t *testing.T) {
	var full, prefix string

	full = "/hello"
	prefix = "/"
	if !hasPrefixDirAccurate(full, prefix, '/') {
		t.Error(full, prefix)
	}

	full = "/a/b/c"
	prefix = "/a/b"
	if !hasPrefixDirAccurate(full, prefix, '/') {
		t.Error(full, prefix)
	}

	full = "/a/b/c"
	prefix = "/a/b/"
	if !hasPrefixDirAccurate(full, prefix, '/') {
		t.Error(full, prefix)
	}

	full = "/a/b/c"
	prefix = "/a/b/c"
	if !hasPrefixDirAccurate(full, prefix, '/') {
		t.Error(full, prefix)
	}

	full = "/a/b/c"
	prefix = "/a/e"
	if hasPrefixDirAccurate(full, prefix, '/') {
		t.Error(full, prefix)
	}

	full = "/a/b/c"
	prefix = "/a/e/"
	if hasPrefixDirAccurate(full, prefix, '/') {
		t.Error(full, prefix)
	}

	full = "/a/b/cd"
	prefix = "/a/b/c"
	if hasPrefixDirAccurate(full, prefix, '/') {
		t.Error(full, prefix)
	}

	full = "/a/b/c"
	prefix = "/a/b/c/"
	if !hasPrefixDirAccurate(full, prefix, '/') {
		t.Error(full, prefix)
	}

	full = "/a/b/c"
	prefix = "/a/b/d"
	if hasPrefixDirAccurate(full, prefix, '/') {
		t.Error(full, prefix)
	}

	full = "/a/b/c"
	prefix = "/a/b/de"
	if hasPrefixDirAccurate(full, prefix, '/') {
		t.Error(full, prefix)
	}
}

func TestHasPrefixDirNoCase(t *testing.T) {
	var full, prefix string

	full = "/hello"
	prefix = "/"
	if !hasPrefixDirNoCase(full, prefix, '/') {
		t.Error(full, prefix)
	}

	full = "/a/b/c"
	prefix = "/a/b"
	if !hasPrefixDirNoCase(full, prefix, '/') {
		t.Error(full, prefix)
	}

	full = "/a/b/c"
	prefix = "/A/b"
	if !hasPrefixDirNoCase(full, prefix, '/') {
		t.Error(full, prefix)
	}

	full = "/a/b/c"
	prefix = "/a/b/"
	if !hasPrefixDirNoCase(full, prefix, '/') {
		t.Error(full, prefix)
	}

	full = "/a/b/c"
	prefix = "/a/B/"
	if !hasPrefixDirNoCase(full, prefix, '/') {
		t.Error(full, prefix)
	}

	full = "/a/b/c"
	prefix = "/a/b/c"
	if !hasPrefixDirNoCase(full, prefix, '/') {
		t.Error(full, prefix)
	}

	full = "/a/b/c"
	prefix = "/a/b/C"
	if !hasPrefixDirNoCase(full, prefix, '/') {
		t.Error(full, prefix)
	}

	full = "/a/b/c"
	prefix = "/a/e"
	if hasPrefixDirNoCase(full, prefix, '/') {
		t.Error(full, prefix)
	}

	full = "/a/b/c"
	prefix = "/A/e"
	if hasPrefixDirNoCase(full, prefix, '/') {
		t.Error(full, prefix)
	}

	full = "/a/b/c"
	prefix = "/a/e/"
	if hasPrefixDirNoCase(full, prefix, '/') {
		t.Error(full, prefix)
	}

	full = "/a/b/c"
	prefix = "/A/E/"
	if hasPrefixDirNoCase(full, prefix, '/') {
		t.Error(full, prefix)
	}

	full = "/a/b/cd"
	prefix = "/a/b/c"
	if hasPrefixDirNoCase(full, prefix, '/') {
		t.Error(full, prefix)
	}

	full = "/a/b/cd"
	prefix = "/a/b/C"
	if hasPrefixDirNoCase(full, prefix, '/') {
		t.Error(full, prefix)
	}

	full = "/a/b/c"
	prefix = "/a/b/c/"
	if !hasPrefixDirNoCase(full, prefix, '/') {
		t.Error(full, prefix)
	}

	full = "/a/b/c"
	prefix = "/A/b/C/"
	if !hasPrefixDirNoCase(full, prefix, '/') {
		t.Error(full, prefix)
	}

	full = "/a/b/c"
	prefix = "/a/b/d"
	if hasPrefixDirNoCase(full, prefix, '/') {
		t.Error(full, prefix)
	}

	full = "/a/b/c"
	prefix = "/a/B/d"
	if hasPrefixDirNoCase(full, prefix, '/') {
		t.Error(full, prefix)
	}

	full = "/a/b/c"
	prefix = "/a/b/de"
	if hasPrefixDirNoCase(full, prefix, '/') {
		t.Error(full, prefix)
	}

	full = "/a/b/c"
	prefix = "/a/B/DE"
	if hasPrefixDirNoCase(full, prefix, '/') {
		t.Error(full, prefix)
	}
}
