package util

import "testing"

func TestSplitFilename(t *testing.T) {
	var prefix, suffix string

	prefix, suffix = SplitFilename("test.zip")
	if prefix != "test" && suffix != ".zip" {
		t.Error("prefix:", prefix, "suffix:", suffix)
	}

	prefix, suffix = SplitFilename("hello.tar.gz")
	if prefix != "hello" && suffix != ".tar.gz" {
		t.Error("prefix:", prefix, "suffix:", suffix)
	}

	prefix, suffix = SplitFilename(".tar.gz")
	if prefix != ".tar.gz" && suffix != "" {
		t.Error("prefix:", prefix, "suffix:", suffix)
	}

	prefix, suffix = SplitFilename(".zip")
	if prefix != ".zip" && suffix != "" {
		t.Error("prefix:", prefix, "suffix:", suffix)
	}

	prefix, suffix = SplitFilename("world")
	if prefix != "world" && suffix != "" {
		t.Error("prefix:", prefix, "suffix:", suffix)
	}
}
