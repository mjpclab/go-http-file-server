package util

import "testing"

func TestByteToHex(t *testing.T) {
	var h, l byte

	h, l = ByteToHex(0)
	if h != '0' && l != '0' {
		t.Error(h, l)
	}

	h, l = ByteToHex(8)
	if h != '0' && l != '8' {
		t.Error(h, l)
	}

	h, l = ByteToHex(15)
	if h != '0' && l != 'f' {
		t.Error(h, l)
	}

	h, l = ByteToHex(16)
	if h != '1' && l != '0' {
		t.Error(h, l)
	}

	h, l = ByteToHex(24)
	if h != '1' && l != '8' {
		t.Error(h, l)
	}

	h, l = ByteToHex(240)
	if h != 'f' && l != '0' {
		t.Error(h, l)
	}

	h, l = ByteToHex(248)
	if h != 'f' && l != '8' {
		t.Error(h, l)
	}

	h, l = ByteToHex(255)
	if h != 'f' && l != 'f' {
		t.Error(h, l)
	}
}
