package goVirtualHost

func (listeners listeners) find(proto, addr string) *listener {
	for _, l := range listeners {
		if l.proto == proto && l.addr == addr {
			return l
		}
	}
	return nil
}
