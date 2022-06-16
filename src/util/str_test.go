package util

import (
	"bytes"
	"testing"
)

func TestReplaceControllingRune(t *testing.T) {
	var str string
	buf := make([]byte, 0, 64)

	str = "abcdefg"
	buf = EscapeControllingRune(str)
	if !bytes.Equal(buf, []byte(str)) {
		t.Error(string(buf))
	}

	str = "abc\tdef"
	buf = EscapeControllingRune(str)
	if !bytes.Equal(buf, []byte("abc\\tdef")) {
		t.Error(string(buf))
	}

	str = "<\000\a\b\f\n\r\t\v>"
	buf = EscapeControllingRune(str)
	if !bytes.Equal(buf, []byte("<\\0\\a\\b\\f\\n\\r\\t\\v>")) {
		t.Error(string(buf))
	}

	str = string([]byte{'[', 0x0e, 127, ']'})
	buf = EscapeControllingRune(str)
	if !bytes.Equal(buf, []byte("[\\x0e\\x7f]")) {
		t.Error(string(buf))
	}
}
