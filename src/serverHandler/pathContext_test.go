package serverHandler

import (
	"testing"
)

func TestPathContext(t *testing.T) {
	var result string

	ctx := &pathContext{}
	result = ctx.QueryString()
	if result != "" {
		t.Error(result)
	}

	ctx.sort = "/n"
	result = ctx.QueryString()
	if result != "?sort=/n" {
		t.Error(result)
	}

}
