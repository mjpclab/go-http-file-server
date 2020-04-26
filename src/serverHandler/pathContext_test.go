package serverHandler

import (
	"testing"
)

func TestPathContext(t *testing.T) {
	var result string
	var sort string

	result = (&pathContext{}).QueryString()
	if result != "" {
		t.Error(result)
	}

	result = (&pathContext{defaultSort: "/n"}).QueryString()
	if result != "" {
		t.Error(result)
	}

	sort = ""
	result = (&pathContext{defaultSort: "/n", sort: &sort}).QueryString()
	if result != "?sort=" {
		t.Error(result)
	}

	sort = "/n"
	result = (&pathContext{defaultSort: "/n", sort: &sort}).QueryString()
	if result != "" {
		t.Error(result)
	}

	sort = ""
	result = (&pathContext{sort: &sort}).QueryString()
	if result != "" {
		t.Error(result)
	}

	sort = "/n"
	result = (&pathContext{sort: &sort}).QueryString()
	if result != "?sort=/n" {
		t.Error(result)
	}
}
