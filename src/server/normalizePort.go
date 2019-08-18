package server

import "strings"

func getDefaultPort(useTLS bool) string {
	if useTLS {
		return ":443"
	} else {
		return ":80"
	}
}

func normalizePort(listen string, useTLS bool) string {
	if len(listen) == 0 {
		return getDefaultPort(useTLS)
	}

	if listen[0] == ':' {
		return listen
	}

	squareEnd := strings.IndexByte(listen, ']')
	isIPv6 := listen[0] == '[' && squareEnd > 0
	isIPv4 := !isIPv6 && strings.IndexByte(listen, '.') >= 0

	if isIPv6 {
		if strings.LastIndexByte(listen, ':') > squareEnd {
			return listen
		}
		return listen + getDefaultPort(useTLS)
	} else if isIPv4 {
		if strings.IndexByte(listen, ':') >= 0 {
			return listen
		}
		return listen + getDefaultPort(useTLS)
	}

	return ":" + listen
}
