package goVirtualHost

import (
	"errors"
	"net"
)

var unknownIPVersion = errors.New("unknown IP version")

func isPrivateIPv4(netIPv4 net.IP) bool {
	return netIPv4[0] == 10 ||
		(netIPv4[0] == 172 && netIPv4[1]&0xf0 == 16) ||
		(netIPv4[0] == 192 && netIPv4[1] == 168)
}

func isPrivateIPv6(netIPv6 net.IP) bool {
	return netIPv6[0]&0xfe == 0xfc
}

func newIPAddr(netIP net.IP) (*ipAddr, error) {
	var version int
	var isNonPrivate bool

	if netIPv4 := netIP.To4(); netIPv4 != nil {
		version = ip4ver
		isNonPrivate = !isPrivateIPv4(netIPv4)
	} else if netIPv6 := netIP.To16(); netIPv6 != nil {
		version = ip6ver
		isNonPrivate = !isPrivateIPv6(netIPv6)
	} else {
		return nil, unknownIPVersion
	}

	instance := &ipAddr{
		netIP:              netIP,
		version:            version,
		isGlobalUnicast:    netIP.IsGlobalUnicast(),
		isLinkLocalUnicast: netIP.IsLinkLocalUnicast(),
		isNonPrivate:       isNonPrivate,
		isNonLoopback:      !netIP.IsLoopback(),
	}
	return instance, nil
}

func (addr *ipAddr) String() string {
	if addr.version == ip6ver {
		return "[" + addr.netIP.String() + "]"
	} else {
		return addr.netIP.String()
	}
}
