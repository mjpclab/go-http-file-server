package util

import "testing"

func TestIsTruthyEnvValue(t *testing.T) {
	expect := func(input string, expectResult bool) bool {
		result := isTruthyEnvValue(input)
		return result == expectResult
	}

	if !expect("", false) {
		t.Error()
	}

	if !expect("false", false) {
		t.Error()
	}

	if !expect("fAlse", false) {
		t.Error()
	}

	if !expect("true", true) {
		t.Error()
	}

	if !expect("truE", true) {
		t.Error()
	}

	if !expect("1", true) {
		t.Error()
	}

	if !expect("0", false) {
		t.Error()
	}

	if !expect("0000", false) {
		t.Error()
	}

	if !expect(" ", false) {
		t.Error()
	}

	if !expect(" 0", false) {
		t.Error()
	}
}
