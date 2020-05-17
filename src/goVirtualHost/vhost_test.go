package goVirtualHost

import "testing"

func TestMatchHostName(t *testing.T) {

	var vh *vhost

	vh = newVhost(nil, []string{"www.example.com"}, nil)
	if !vh.matchHostName("www.example.com") {
		t.Error()
	}
	if vh.matchHostName("example.com") {
		t.Error()
	}

	vh = newVhost(nil, []string{".example.com"}, nil)
	if !vh.matchHostName("www.example.com") {
		t.Error()
	}
	if vh.matchHostName("example.com") {
		t.Error()
	}

	vh = newVhost(nil, []string{".example.com", "example.com"}, nil)
	if !vh.matchHostName("www.example.com") {
		t.Error()
	}
	if !vh.matchHostName("example.com") {
		t.Error()
	}

	vh = newVhost(nil, []string{"example."}, nil)
	if !vh.matchHostName("example.com") {
		t.Error()
	}
	if !vh.matchHostName("example.net") {
		t.Error()
	}
}
