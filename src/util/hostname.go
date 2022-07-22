package util

import (
	"../shimgo"
	"strings"
)

func ExtractHostnamePort(host string) (hostname, port string) {
	if len(host) == 0 {
		return
	}

	hostname = host
	if hostname[len(hostname)-1] != ']' { // not [IPv6]
		if colonIndex := strings.LastIndex(hostname, "]:"); colonIndex > 0 { // [IPv6]:port
			hostname = hostname[:colonIndex+1]
		} else if colonIndex = shimgo.Strings_LastIndexByte(hostname, ':'); colonIndex >= 0 {
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
	} else if colonIndex = shimgo.Strings_LastIndexByte(listen, ':'); colonIndex >= 0 {
		return listen[colonIndex+1:]
	} else {
		lenListen := len(listen)
		if (lenListen < 5 && IsDigits(listen)) ||
			(lenListen == 5 && listen < "65536") {
			return listen
		}
	}

	return ""
}

const protoSuffix string = "://"
const protoSuffixLen = len(protoSuffix)

func ExtractHostFromUrl(url string) string {
	protoIndex := strings.Index(url, protoSuffix)
	if protoIndex >= 0 {
		url = url[protoIndex+protoSuffixLen:]
	}

	slashIndex := strings.IndexByte(url, '/')
	if slashIndex >= 0 {
		url = url[:slashIndex]
	}

	url = strings.ToLower(url)

	return url
}

func ExtractHostsFromUrls(urls []string) []string {
	hosts := make([]string, 0, len(urls))

	for i := range urls {
		host := ExtractHostFromUrl(urls[i])
		if len(host) > 0 {
			hosts = append(hosts, host)
		}
	}

	return hosts
}
