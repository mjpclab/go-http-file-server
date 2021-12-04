package goVirtualHost

import (
	"net"
	"sort"
	"strings"
)

var whitespaceRemover = strings.NewReplacer(
	" ", "",
	"\t", "",
	"\v", "",
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

func isDefaultPort(port string, useTLS bool) bool {
	if useTLS {
		return port == ":443"
	} else {
		return port == ":80"
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

func splitListen(listen string, useTLS bool) (l43proto, ip, port string) {
	listen = whitespaceRemover.Replace(listen)

	// empty, use default tcp port
	if len(listen) == 0 {
		return tcp46, "", getDefaultPort(useTLS)
	}

	// :port
	if listen[0] == ':' {
		return tcp46, "", listen
	}

	// port
	if isDigits(listen) {
		return tcp46, "", ":" + listen
	}

	// unix socket path
	if strings.IndexByte(listen, '/') >= 0 {
		return unix, "", listen
	}

	colonIndex := strings.IndexByte(listen, ':')
	lastColonIndex := strings.LastIndexByte(listen, ':')

	// ipv6
	squareEnd := strings.IndexByte(listen, ']')
	isIPv6 := listen[0] == '[' && squareEnd > 0 && colonIndex > 0 && colonIndex < squareEnd
	if isIPv6 {
		var ip, port string
		if lastColonIndex == squareEnd+1 { // has port number
			ip = listen[:lastColonIndex]
			port = listen[lastColonIndex:]
		} else {
			ip = listen
			port = getDefaultPort(useTLS)
		}
		if isWildcardIPv6(ip) {
			ip = ""
		}
		return tcp6, ip, port
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
		var ip, port string
		if colonIndex >= 0 { // has port number
			ip = listen[:colonIndex]
			port = listen[colonIndex:]
		} else {
			ip = listen
			port = getDefaultPort(useTLS)
		}
		if isWildcardIPv4(ip) {
			ip = ""
		}
		return tcp4, ip, port
	}

	// suppose to be a domain with port
	if colonIndex >= 0 {
		return tcp46, listen[:colonIndex], listen[colonIndex:]
	}

	// suppose to be a domain
	return tcp46, listen, getDefaultPort(useTLS)
}

func isWildcardIPv4(ip string) bool {
	return ip == "0.0.0.0"
}

func isWildcardIPv6(ip string) bool {
	// min len==4, [::]
	// max len==41, [0000:0000:0000:0000:0000:0000:0000:0000]
	if len(ip) < 4 || len(ip) > 41 {
		return false
	}

	// remove brackets
	ip = ip[1 : len(ip)-1]

	for _, c := range ip {
		switch c {
		case '0', ':':
			continue
		}
		return false
	}

	return true
}

func getAllIfaceIPs(includeLoopback bool) (all, allv4, allv6 []string) {
	var allAddrs, allAddrsV4, allAddrsV6 ipAddrs

	netAddrs, _ := net.InterfaceAddrs()
	for _, netAddr := range netAddrs {
		var netIP net.IP
		switch v := netAddr.(type) {
		case *net.IPNet:
			netIP = v.IP
		case *net.IPAddr:
			netIP = v.IP
		default:
			continue
		}

		addr, _ := newIPAddr(netIP)
		if addr.isNonLoopback || includeLoopback {
			allAddrs = append(allAddrs, addr)
			if addr.version == ip4ver {
				allAddrsV4 = append(allAddrsV4, addr)
			} else if addr.version == ip6ver {
				allAddrsV6 = append(allAddrsV6, addr)
			}
		}
	}

	sort.Sort(allAddrs)
	all = allAddrs.String()

	sort.Sort(allAddrsV4)
	allv4 = allAddrsV4.String()

	sort.Sort(allAddrsV6)
	allv6 = allAddrsV6.String()
	return
}
