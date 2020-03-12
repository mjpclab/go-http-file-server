package util

import "testing"

func TestCompareNumInFilename(t *testing.T) {
	var prev, next string

	prev = "2"
	next = "3"
	if !CompareNumInFilename([]byte(prev), []byte(next)) {
		t.Error(prev, next)
	}

	prev = "2"
	next = "20"
	if !CompareNumInFilename([]byte(prev), []byte(next)) {
		t.Error(prev, next)
	}

	prev = "1.1"
	next = "1.10"
	if !CompareNumInFilename([]byte(prev), []byte(next)) {
		t.Error(prev, next)
	}

	prev = "1.2"
	next = "1.10"
	if !CompareNumInFilename([]byte(prev), []byte(next)) {
		t.Error(prev, next)
	}

	prev = "1.3.2"
	next = "1.3.10"
	if !CompareNumInFilename([]byte(prev), []byte(next)) {
		t.Error(prev, next)
	}

	prev = "name"
	next = "name-suffix"
	if !CompareNumInFilename([]byte(prev), []byte(next)) {
		t.Error(prev, next)
	}

	prev = "name.txt"
	next = "name-1.txt"
	if !CompareNumInFilename([]byte(prev), []byte(next)) {
		t.Error(prev, next)
	}
}

func TestExtractPrefixDigits(t *testing.T) {
	result := string(extractPrefixDigits([]byte("123.haha")))
	if result != "123" {
		t.Error(result)
	}
}
