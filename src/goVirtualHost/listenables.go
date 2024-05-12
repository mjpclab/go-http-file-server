package goVirtualHost

func (ls listenables) find(proto, ip, port string) *listenable {
	for _, l := range ls {
		if l.proto == proto && l.ip == ip && l.port == port {
			return l
		}
	}
	return nil
}
