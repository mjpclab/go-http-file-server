package util

import "testing"

func TestExtractHostnamePort(t *testing.T) {
	var hostname, port string

	hostname, port = ExtractHostnamePort("www.example.com:123")
	if hostname != "www.example.com" || port != ":123" {
		t.Error(hostname, port)
	}

	hostname, port = ExtractHostnamePort("www.example.net")
	if hostname != "www.example.net" || port != "" {
		t.Error(hostname, port)
	}

	hostname, port = ExtractHostnamePort(":234")
	if hostname != "" || port != ":234" {
		t.Error(hostname, port)
	}

	hostname, port = ExtractHostnamePort("[fe80::1]:345")
	if hostname != "[fe80::1]" || port != ":345" {
		t.Error(hostname, port)
	}

	hostname, port = ExtractHostnamePort("[fe80::1]")
	if hostname != "[fe80::1]" || port != "" {
		t.Error(hostname, port)
	}
}

func TestExtractListenPort(t *testing.T) {
	if ExtractListenPort("abc:123") != "123" {
		t.Error(1)
	}

	if ExtractListenPort("127.0.0.1:123") != "123" {
		t.Error(1)
	}

	if ExtractListenPort("[::1]:123") != "123" {
		t.Error(1)
	}

	if ExtractListenPort("abc") != "" {
		t.Error(1)
	}

	if ExtractListenPort("127.0.0.1") != "" {
		t.Error(1)
	}

	if ExtractListenPort("[::1]") != "" {
		t.Error(1)
	}

	if ExtractListenPort(":123") != "123" {
		t.Error(1)
	}

	if ExtractListenPort("123") != "123" {
		t.Error(1)
	}

	if ExtractListenPort("65535") != "65535" {
		t.Error(1)
	}

	if ExtractListenPort("65536") != "" {
		t.Error(1)
	}
}
