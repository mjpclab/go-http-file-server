package util

import "testing"

func TestIsDigits(t *testing.T) {
	var str string
	str = "12345"
	if !IsDigits(str) {
		t.Fail()
	}

	str = "x12345"
	if IsDigits(str) {
		t.Fail()
	}
}
