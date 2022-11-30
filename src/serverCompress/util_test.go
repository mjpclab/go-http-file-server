package serverCompress

import "testing"

func TestIsCompressibleType(t *testing.T) {
	if !isCompressibleType("text/plain") {
		t.Error()
	}
	if !isCompressibleType("text/plain; charset=utf-8") {
		t.Error()
	}

	if !isCompressibleType("image/svg+xml") {
		t.Error()
	}
	if !isCompressibleType("image/svg+xml; charset=utf-8") {
		t.Error()
	}
	if isCompressibleType("image/png") {
		t.Error()
	}
	if isCompressibleType("image/png; foo=bar") {
		t.Error()
	}

	if !isCompressibleType("application/javascript") {
		t.Error()
	}
	if !isCompressibleType("application/javascript; charset=utf-8") {
		t.Error()
	}
}
