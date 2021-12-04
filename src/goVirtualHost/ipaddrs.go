package goVirtualHost

func (addrs ipAddrs) Len() int {
	return len(addrs)
}

func (addrs ipAddrs) Less(i, j int) bool {
	addr1 := addrs[i]
	addr2 := addrs[j]

	if addr1.isGlobalUnicast != addr2.isGlobalUnicast {
		return addr1.isGlobalUnicast
	}

	if addr1.isLinkLocalUnicast != addr2.isLinkLocalUnicast {
		return addr1.isLinkLocalUnicast
	}

	if addr1.isNonPrivate != addr2.isNonPrivate {
		return addr1.isNonPrivate
	}

	if addr1.isNonLoopback != addr2.isNonLoopback {
		return addr1.isNonLoopback
	}

	if addr1.version != addr2.version {
		return addr1.version == ip4ver
	}

	return i < j
}

func (addrs ipAddrs) Swap(i, j int) {
	addrs[i], addrs[j] = addrs[j], addrs[i]
}

func (addrs ipAddrs) String() []string {
	results := make([]string, len(addrs))
	for i := range addrs {
		results[i] = addrs[i].String()
	}
	return results
}
