package util

import "testing"

func TestCompareNumInFilename(t *testing.T) {
	var prev, next string
	var less bool

	prev = "2"
	next = "3"
	less, _ = CompareNumInFilename([]byte(prev), []byte(next))
	if !less {
		t.Error(prev, next)
	}

	prev = "2"
	next = "20"
	less, _ = CompareNumInFilename([]byte(prev), []byte(next))
	if !less {
		t.Error(prev, next)
	}

	prev = "1.1"
	next = "1.10"
	less, _ = CompareNumInFilename([]byte(prev), []byte(next))
	if !less {
		t.Error(prev, next)
	}

	prev = "1.2"
	next = "1.10"
	less, _ = CompareNumInFilename([]byte(prev), []byte(next))
	if !less {
		t.Error(prev, next)
	}

	prev = "1.3.2"
	next = "1.3.10"
	less, _ = CompareNumInFilename([]byte(prev), []byte(next))
	if !less {
		t.Error(prev, next)
	}

	prev = "name"
	next = "name-suffix"
	less, _ = CompareNumInFilename([]byte(prev), []byte(next))
	if !less {
		t.Error(prev, next)
	}

	prev = "name.txt"
	next = "name-1.txt"
	less, _ = CompareNumInFilename([]byte(prev), []byte(next))
	if !less {
		t.Error(prev, next)
	}

	prev = "a1.txt"
	next = "B1.txt"
	less, _ = CompareNumInFilename([]byte(prev), []byte(next))
	if !less {
		t.Error(prev, next)
	}
}

func TestExtractPrefixDigits(t *testing.T) {
	result, _ := extractPrefixDigits([]byte("123.haha"))
	strResult := string(result)
	if strResult != "123" {
		t.Error(result)
	}
}
