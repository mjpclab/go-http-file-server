package util

import "strings"

func ExtractHostname(host string) string {
	colonIndex := strings.LastIndexByte(host, ':')
	if colonIndex >= 0 {
		return host[:colonIndex]
	}
	return host
}
