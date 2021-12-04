package goVirtualHost

import "testing"

func TestExtractHostname(t *testing.T) {
	var host, hostname string

	host = "example.com"
	hostname = extractHostName(host)
	if hostname != "example.com" {
		t.Error(hostname)
	}

	host = "example.com:8080"
	hostname = extractHostName(host)
	if hostname != "example.com" {
		t.Error(hostname)
	}

	host = "[fe80::1]"
	hostname = extractHostName(host)
	if hostname != "fe80::1" {
		t.Error(hostname)
	}

	host = "[fe80::1]:8080"
	hostname = extractHostName(host)
	if hostname != "fe80::1" {
		t.Error(hostname)
	}
}

func TestNormalizeHostNames(t *testing.T) {
	inputs := []string{"aA", "", "Bb"}
	results := normalizeHostNames(inputs)
	if len(results) != 2 || results[0] != "aa" || results[1] != "bb" {
		t.Error(results)
	}
}

func TestSplitListen(t *testing.T) {
	match := func(listen string, useTLS bool, expectProto, expectIP, expectPort string) bool {
		proto, ip, port := splitListen(listen, useTLS)
		return proto == expectProto && ip == expectIP && port == expectPort
	}

	// ipv4
	if !match("1.2.3.4", false, "tcp4", "1.2.3.4", ":80") {
		t.Error()
	}
	if !match(" 1.2.3.4\t", false, "tcp4", "1.2.3.4", ":80") {
		t.Error()
	}
	if !match("1.2.3.4", true, "tcp4", "1.2.3.4", ":443") {
		t.Error()
	}
	if !match("0.0.0.0", false, "tcp4", "", ":80") {
		t.Error()
	}

	// ipv4:port
	if !match("2.3.4.5:6", false, "tcp4", "2.3.4.5", ":6") {
		t.Error()
	}
	if !match("2.3.4.5:6", true, "tcp4", "2.3.4.5", ":6") {
		t.Error()
	}
	if !match("0.0.0.0:6", false, "tcp4", "", ":6") {
		t.Error()
	}

	// ipv6
	if !match("[::1]", false, "tcp6", "[::1]", ":80") {
		t.Error()
	}
	if !match("[::1]", true, "tcp6", "[::1]", ":443") {
		t.Error()
	}
	if !match("[::]", false, "tcp6", "", ":80") {
		t.Error()
	}

	// ipv6:port
	if !match("[fe80::1234]:7", false, "tcp6", "[fe80::1234]", ":7") {
		t.Error()
	}
	if !match("[fe80::1234]:7", true, "tcp6", "[fe80::1234]", ":7") {
		t.Error()
	}
	if !match("[::]:7", false, "tcp6", "", ":7") {
		t.Error()
	}

	// port
	if !match("8080", false, "tcp", "", ":8080") {
		t.Error()
	}
	if !match("8080", true, "tcp", "", ":8080") {
		t.Error()
	}

	// :port
	if !match(":3000", false, "tcp", "", ":3000") {
		t.Error()
	}
	if !match(":3000", true, "tcp", "", ":3000") {
		t.Error()
	}

	// hostname
	if !match("example.com", false, "tcp", "example.com", ":80") {
		t.Error()
	}
	if !match("example.com", true, "tcp", "example.com", ":443") {
		t.Error()
	}

	// hostname:port
	if !match("example.com:3210", false, "tcp", "example.com", ":3210") {
		t.Error()
	}
	if !match("example.com:3210", true, "tcp", "example.com", ":3210") {
		t.Error()
	}

	// socket
	if !match("/var/run/ghfs.sock", false, "unix", "", "/var/run/ghfs.sock") {
		t.Error()
	}
	if !match("/var/run/ghfs.sock", true, "unix", "", "/var/run/ghfs.sock") {
		t.Error()
	}
}

func TestIsIPv4Wildcard(t *testing.T) {
	if !isWildcardIPv4("0.0.0.0") {
		t.Error()
	}
}

func TestIsIPv6Wildcard(t *testing.T) {
	if isWildcardIPv6("0.0.0.0") {
		t.Error()
	}

	if isWildcardIPv6("[fe80::1]") {
		t.Error()
	}

	if isWildcardIPv6("[::1]") {
		t.Error()
	}

	if isWildcardIPv6("[fe80::1]:8080") {
		t.Error()
	}

	if isWildcardIPv6("[::1]:8080") {
		t.Error()
	}

	if isWildcardIPv6("8080") {
		t.Error()
	}

	if isWildcardIPv6(":8080") {
		t.Error()
	}

	if isWildcardIPv6("[::]:8080") {
		t.Error()
	}

	if isWildcardIPv6("[::0]:8080") {
		t.Error()
	}

	if !isWildcardIPv6("[::]") {
		t.Error()
	}

	if !isWildcardIPv6("[::0]") {
		t.Error()
	}

	if !isWildcardIPv6("[0::]") {
		t.Error()
	}

	if !isWildcardIPv6("[0::0]") {
		t.Error()
	}

	if !isWildcardIPv6("[::00]") {
		t.Error()
	}

	if !isWildcardIPv6("[00::00]") {
		t.Error()
	}

	if !isWildcardIPv6("[0000:0000:0000:0000:0000:0000:0000:0000]") {
		t.Error()
	}
}
