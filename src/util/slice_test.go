package util

import "testing"

func TestContains(t *testing.T) {
	coll := []string{"abc", "def", "ghi"}
	found := Contains(coll, "def")
	if !found {
		t.Error()
	}

	found = Contains(coll, "DEF")
	if found {
		t.Error()
	}
}
