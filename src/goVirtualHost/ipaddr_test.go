package goVirtualHost

import (
	"net"
	"testing"
)

func TestIsPrivateIPv4(t *testing.T) {
	var ip net.IP

	ip = net.IPv4(10, 1, 2, 3).To4()
	if !isPrivateIPv4(ip) {
		t.Error()
	}

	ip = net.IPv4(172, 16, 2, 3).To4()
	if !isPrivateIPv4(ip) {
		t.Error()
	}

	ip = net.IPv4(172, 17, 2, 3).To4()
	if !isPrivateIPv4(ip) {
		t.Error()
	}

	ip = net.IPv4(172, 30, 2, 3).To4()
	if !isPrivateIPv4(ip) {
		t.Error()
	}

	ip = net.IPv4(172, 31, 2, 3).To4()
	if !isPrivateIPv4(ip) {
		t.Error()
	}

	ip = net.IPv4(192, 168, 4, 5).To4()
	if !isPrivateIPv4(ip) {
		t.Error()
	}

	ip = net.IPv4(8, 8, 8, 8).To4()
	if isPrivateIPv4(ip) {
		t.Error()
	}
}

func TestIsPrivateIPv6(t *testing.T) {
	var ip net.IP

	ip = net.IP{0xfc, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
	if !isPrivateIPv6(ip) {
		t.Error()
	}

	ip = net.IP{0xfe, 0x80, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
	if isPrivateIPv6(ip) {
		t.Error()
	}
}
