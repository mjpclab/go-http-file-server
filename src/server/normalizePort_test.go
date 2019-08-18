package server

import "testing"

func TestNormalizePort(t *testing.T) {
	// ipv4 only
	ipv4 := "1.2.3.4"

	ipv4Http := normalizePort(ipv4, false)
	if ipv4Http != "1.2.3.4:80" {
		t.Error(ipv4Http)
	}

	ipv4Https := normalizePort(ipv4, true)
	if ipv4Https != "1.2.3.4:443" {
		t.Error(ipv4Https)
	}

	// ipv4 with port
	ipv4Port := "2.3.4.5:6"

	ipv4PortHttp := normalizePort(ipv4Port, false)
	if ipv4PortHttp != ipv4Port {
		t.Error(ipv4PortHttp)
	}

	ipv4PortHttps := normalizePort(ipv4Port, true)
	if ipv4PortHttps != ipv4Port {
		t.Error(ipv4PortHttps)
	}

	// ipv6 only
	ipv6 := "[::1]"

	ipv6Http := normalizePort(ipv6, false)
	if ipv6Http != "[::1]:80" {
		t.Error(ipv6Http)
	}

	ipv6Https := normalizePort(ipv6, true)
	if ipv6Https != "[::1]:443" {
		t.Error(ipv6Https)
	}

	// ipv6 with port
	ipv6Port := "[fe80::1234]:7"

	ipv6PortHttp := normalizePort(ipv6Port, false)
	if ipv6PortHttp != ipv6Port {
		t.Error(ipv6PortHttp)
	}

	ipv6PortHttps := normalizePort(ipv6Port, true)
	if ipv6PortHttps != ipv6Port {
		t.Error(ipv6PortHttps)
	}

	// port number only
	portNum := "8080"
	portNumHttp := normalizePort(portNum, false)
	if portNumHttp != ":8080" {
		t.Error(portNumHttp)
	}

	portNumHttps := normalizePort(portNum, true)
	if portNumHttps != ":8080" {
		t.Error(portNumHttps)
	}

	// :port
	port := ":3000"
	portHttp := normalizePort(port, false)
	if portHttp != port {
		t.Error(portHttp)
	}

	portHttps := normalizePort(port, true)
	if portHttps != port {
		t.Error(portHttps)
	}
}
