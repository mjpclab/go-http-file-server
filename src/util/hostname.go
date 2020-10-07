package util

import "strings"

func ExtractHostnamePort(host string) (hostname, port string) {
	if len(host) == 0 {
		return
	}

	hostname = host
	if hostname[len(hostname)-1] != ']' { // not [IPv6]
		if colonIndex := strings.LastIndex(hostname, "]:"); colonIndex > 0 { // [IPv6]:port
			hostname = hostname[:colonIndex+1]
		} else if colonIndex = strings.LastIndexByte(hostname, ':'); colonIndex >= 0 {
			hostname = hostname[:colonIndex]
		}
	}

	port = host[len(hostname):]

	return
}

func ExtractListenPort(listen string) string {
	if len(listen) == 0 {
		return ""
	}

	if colonIndex := strings.LastIndex(listen, "]:"); colonIndex > 0 { // [IPv6]:port
		return listen[colonIndex+2:]
	} else if listen[len(listen)-1] == ']' { // [IPv6]
		return ""
	} else if colonIndex = strings.LastIndexByte(listen, ':'); colonIndex >= 0 {
		return listen[colonIndex+1:]
	} else {
		lenListen := len(listen)
		if (lenListen < 5 && IsDigits(listen)) ||
			(lenListen == 5 && strings.Compare(listen, "65535") == -1) {
			return listen
		}
	}

	return ""
}
