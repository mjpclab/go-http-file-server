package goVirtualHost

func (listeners listeners) find(proto, ip, port string) *listener {
	for _, l := range listeners {
		if l.proto == proto && l.ip == ip && l.port == port {
			return l
		}
	}
	return nil
}
