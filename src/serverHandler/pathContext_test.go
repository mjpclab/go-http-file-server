package serverHandler

import (
	"testing"
)

func TestPathContext(t *testing.T) {
	var result string
	var sort string
	var download bool

	result = (&pathContext{}).QueryString()
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

	download = false
	result = (&pathContext{download: download}).QueryString()
	if result != "" {
		t.Error(result)
	}

	download = false
	sort = "/n"
	result = (&pathContext{download: download, sort: &sort}).QueryString()
	if result != "?sort=/n" {
		t.Error(result)
	}

	download = true
	result = (&pathContext{download: download}).QueryString()
	if result != "?download" {
		t.Error(result)
	}

	download = true
	sort = "/n"
	result = (&pathContext{download: download, sort: &sort}).QueryString()
	if result != "?download&sort=/n" {
		t.Error(result)
	}
}
