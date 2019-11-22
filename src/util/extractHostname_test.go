package util

import "testing"

func TestExtractHostname(t *testing.T) {
	var host, hostname string

	host = "example.com"
	hostname = ExtractHostname(host)
	if hostname != "example.com" {
		t.Error(hostname)
	}

	host = "example.com:8080"
	hostname = ExtractHostname(host)
	if hostname != "example.com" {
		t.Error(hostname)
	}
}
