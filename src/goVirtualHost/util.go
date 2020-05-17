package goVirtualHost

import (
	"strings"
)

func extractHostName(host string) string {
	hostLen := len(host)
	if hostLen == 0 {
		return host
	}

	if hostLen >= 5 && host[0] == '[' { // [IPV6]:port, "[" "ip" "]" ":" "port", 5 parts
		maxIndex := hostLen - 1
		closeIndex := strings.IndexByte(host, ']')
		if closeIndex == maxIndex {
			return host[1:closeIndex]
		}
		if closeIndex > 1 && closeIndex < maxIndex && host[closeIndex+1] == ':' {
			return host[1:closeIndex]
		}
	}

	colonIndex := strings.LastIndexByte(host, ':')
	if colonIndex >= 0 {
		return host[:colonIndex]
	}
	return host
}

func normalizeHostNames(inputs []string) []string {
	output := make([]string, 0, len(inputs))

	for _, str := range inputs {
		if len(str) > 0 {
			name := strings.ToLower(str)
			output = append(output, name)
		}
	}

	return output
}

func getDefaultPort(useTLS bool) string {
	if useTLS {
		return ":443"
	} else {
		return ":80"
	}
}

func isDigits(input string) bool {
	for i, length := 0, len(input); i < length; i++ {
		b := input[i]
		if b < '0' || b > '9' {
			return false
		}
	}

	return true
}

func splitListen(listen string, useTLS bool) (proto, addr string) {
	// empty, use default tcp port
	if len(listen) == 0 {
		return "tcp", getDefaultPort(useTLS)
	}

	// :port
	if listen[0] == ':' {
		return "tcp", listen
	}

	// port
	if isDigits(listen) {
		return "tcp", ":" + listen
	}

	// unix socket path
	if strings.IndexByte(listen, '/') >= 0 {
		return "unix", listen
	}

	colonIndex := strings.IndexByte(listen, ':')
	lastColonIndex := strings.LastIndexByte(listen, ':')

	// ipv6
	squareEnd := strings.IndexByte(listen, ']')
	isIPv6 := listen[0] == '[' && squareEnd > 0 && colonIndex > 0 && colonIndex < squareEnd
	if isIPv6 {
		if lastColonIndex > squareEnd { // has port number
			return "tcp6", listen
		}
		return "tcp6", listen + getDefaultPort(useTLS)
	}

	// ipv4
	dotIndex1 := strings.IndexByte(listen, '.')
	dotIndex2 := dotIndex1 + 1 + strings.IndexByte(listen[dotIndex1+1:], '.')
	dotIndex3 := dotIndex2 + 1 + strings.IndexByte(listen[dotIndex2+1:], '.')
	dotIndex4 := dotIndex3 + 1 + strings.IndexByte(listen[dotIndex3+1:], '.')
	lastDotIndex := strings.LastIndexByte(listen, '.')
	isIPv4 := dotIndex1 > 0 && dotIndex2 > dotIndex1 && dotIndex3 > dotIndex2 && dotIndex4 == dotIndex3 &&
		colonIndex == lastColonIndex &&
		(lastColonIndex == -1 || lastColonIndex > lastDotIndex+1)
	if isIPv4 {
		if colonIndex >= 0 { // has port number
			return "tcp4", listen
		}
		return "tcp4", listen + getDefaultPort(useTLS)
	}

	// suppose to be a domain with port
	if colonIndex >= 0 {
		return "tcp", listen
	}

	// suppose to be a domain
	return "tcp", listen + getDefaultPort(useTLS)
}
