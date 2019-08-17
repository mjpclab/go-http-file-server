package fmtSize

import (
	"testing"
)

func TestFmtUnit(t *testing.T) {
	v := 1536
	strV := fmtUnit('K', KB, int64(v))
	if strV != "1K" {
		t.Error("value is not '1K'")
	}
}
