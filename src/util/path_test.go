package util

import "testing"

func TestHasUrlPrefixDir(t *testing.T) {
	var full, prefix string

	full = "/a/b/c"
	prefix = "/a/b"
	if !HasUrlPrefixDir(full, prefix) {
		t.Error(full, prefix)
	}

	full = "/a/b/c"
	prefix = "/a/b/"
	if !HasUrlPrefixDir(full, prefix) {
		t.Error(full, prefix)
	}

	full = "/a/b/c"
	prefix = "/a/b/c"
	if !HasUrlPrefixDir(full, prefix) {
		t.Error(full, prefix)
	}

	full = "/a/b/c"
	prefix = "/a/e"
	if HasUrlPrefixDir(full, prefix) {
		t.Error(full, prefix)
	}

	full = "/a/b/c"
	prefix = "/a/e/"
	if HasUrlPrefixDir(full, prefix) {
		t.Error(full, prefix)
	}

	full = "/a/b/cd"
	prefix = "/a/b/c"
	if HasUrlPrefixDir(full, prefix) {
		t.Error(full, prefix)
	}

	full = "/a/b/c"
	prefix = "/a/b/c/"
	if !HasUrlPrefixDir(full, prefix) {
		t.Error(full, prefix)
	}

	full = "/a/b/c"
	prefix = "/a/b/d"
	if HasUrlPrefixDir(full, prefix) {
		t.Error(full, prefix)
	}

	full = "/a/b/c"
	prefix = "/a/b/de"
	if HasUrlPrefixDir(full, prefix) {
		t.Error(full, prefix)
	}
}

func TestHasUrlPrefixDirNoCase(t *testing.T) {
	var full, prefix string

	full = "/a/b/c"
	prefix = "/a/b"
	if !HasUrlPrefixDirNoCase(full, prefix) {
		t.Error(full, prefix)
	}

	full = "/a/b/c"
	prefix = "/A/b"
	if !HasUrlPrefixDirNoCase(full, prefix) {
		t.Error(full, prefix)
	}

	full = "/a/b/c"
	prefix = "/a/b/"
	if !HasUrlPrefixDirNoCase(full, prefix) {
		t.Error(full, prefix)
	}

	full = "/a/b/c"
	prefix = "/a/B/"
	if !HasUrlPrefixDirNoCase(full, prefix) {
		t.Error(full, prefix)
	}

	full = "/a/b/c"
	prefix = "/a/b/c"
	if !HasUrlPrefixDirNoCase(full, prefix) {
		t.Error(full, prefix)
	}

	full = "/a/b/c"
	prefix = "/a/b/C"
	if !HasUrlPrefixDirNoCase(full, prefix) {
		t.Error(full, prefix)
	}

	full = "/a/b/c"
	prefix = "/a/e"
	if HasUrlPrefixDirNoCase(full, prefix) {
		t.Error(full, prefix)
	}

	full = "/a/b/c"
	prefix = "/A/e"
	if HasUrlPrefixDirNoCase(full, prefix) {
		t.Error(full, prefix)
	}

	full = "/a/b/c"
	prefix = "/a/e/"
	if HasUrlPrefixDirNoCase(full, prefix) {
		t.Error(full, prefix)
	}

	full = "/a/b/c"
	prefix = "/A/E/"
	if HasUrlPrefixDirNoCase(full, prefix) {
		t.Error(full, prefix)
	}

	full = "/a/b/cd"
	prefix = "/a/b/c"
	if HasUrlPrefixDirNoCase(full, prefix) {
		t.Error(full, prefix)
	}

	full = "/a/b/cd"
	prefix = "/a/b/C"
	if HasUrlPrefixDirNoCase(full, prefix) {
		t.Error(full, prefix)
	}

	full = "/a/b/c"
	prefix = "/a/b/c/"
	if !HasUrlPrefixDirNoCase(full, prefix) {
		t.Error(full, prefix)
	}

	full = "/a/b/c"
	prefix = "/A/b/C/"
	if !HasUrlPrefixDirNoCase(full, prefix) {
		t.Error(full, prefix)
	}

	full = "/a/b/c"
	prefix = "/a/b/d"
	if HasUrlPrefixDirNoCase(full, prefix) {
		t.Error(full, prefix)
	}

	full = "/a/b/c"
	prefix = "/a/B/d"
	if HasUrlPrefixDirNoCase(full, prefix) {
		t.Error(full, prefix)
	}

	full = "/a/b/c"
	prefix = "/a/b/de"
	if HasUrlPrefixDirNoCase(full, prefix) {
		t.Error(full, prefix)
	}

	full = "/a/b/c"
	prefix = "/a/B/DE"
	if HasUrlPrefixDirNoCase(full, prefix) {
		t.Error(full, prefix)
	}
}
