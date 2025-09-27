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

	sort = ""
	result = (&pathContext{defaultSort: "/n", sort: sort}).QueryString()
	if result != "" {
		t.Error(result)
	}

	sort = "/n"
	result = (&pathContext{defaultSort: "/n", sort: sort}).QueryString()
	if result != "" {
		t.Error(result)
	}

	sort = ""
	result = (&pathContext{sort: sort}).QueryString()
	if result != "" {
		t.Error(result)
	}

	sort = "/n"
	result = (&pathContext{sort: sort}).QueryString()
	if result != "?sort=/n" {
		t.Error(result)
	}

	result = (&pathContext{simple: false}).QueryString()
	if result != "" {
		t.Error(result)
	}

	result = (&pathContext{simple: true}).QueryString()
	if result != "?simple" {
		t.Error(result)
	}

	result = (&pathContext{download: true}).QueryString()
	if result != "?download" {
		t.Error(result)
	}

	result = (&pathContext{simple: true, download: true}).QueryString()
	if result != "?simple&download" {
		t.Error(result)
	}

	sort = "/n"
	result = (&pathContext{simple: false, sort: sort}).QueryString()
	if result != "?sort=/n" {
		t.Error(result)
	}

	sort = "/n"
	result = (&pathContext{simple: true, sort: sort}).QueryString()
	if result != "?simple&sort=/n" {
		t.Error(result)
	}

	sort = "/n"
	result = (&pathContext{download: true, sort: sort}).QueryString()
	if result != "?download&sort=/n" {
		t.Error(result)
	}

	sort = "/n"
	result = (&pathContext{simple: true, download: true, sort: sort}).QueryString()
	if result != "?simple&download&sort=/n" {
		t.Error(result)
	}
}
