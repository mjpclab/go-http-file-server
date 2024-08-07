package goVirtualHost

import "testing"

func TestMatchHostName(t *testing.T) {

	var vh *vhost

	vh = newVhost([]string{"www.example.com"}, nil, nil, nil)
	if !vh.matchHostName("www.example.com") {
		t.Error()
	}
	if vh.matchHostName("example.com") {
		t.Error()
	}

	vh = newVhost([]string{".example.com"}, nil, nil, nil)
	if !vh.matchHostName("www.example.com") {
		t.Error()
	}
	if vh.matchHostName("example.com") {
		t.Error()
	}

	vh = newVhost([]string{".example.com", "example.com"}, nil, nil, nil)
	if !vh.matchHostName("www.example.com") {
		t.Error()
	}
	if !vh.matchHostName("example.com") {
		t.Error()
	}

	vh = newVhost([]string{"example."}, nil, nil, nil)
	if !vh.matchHostName("example.com") {
		t.Error()
	}
	if !vh.matchHostName("example.net") {
		t.Error()
	}
}
