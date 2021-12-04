package goVirtualHost

import (
	"errors"
	"sync"
)

var alreadyOpened = errors.New("already opened")

func NewService() *Service {
	service := &Service{
		state:     statePrepare,
		listeners: listeners{},
		servers:   servers{},
		vhosts:    vhosts{},
	}

	return service
}

func (svc *Service) addVhostToServers(vhost *vhost, params params) {
	for _, param := range params {
		// listeners, servers
		var listener *listener
		var server *server

		listener = svc.listeners.find(param.proto, param.ip, param.port)
		if listener != nil {
			server = listener.server
		} else {
			server = newServer(param.useTLS)
			listener = newListener(param.proto, param.ip, param.port)
			listener.server = server

			svc.listeners = append(svc.listeners, listener)
			svc.servers = append(svc.servers, server)
		}

		// server -> vhost
		server.vhosts = append(server.vhosts, vhost)
	}
}

func (svc *Service) Add(info *HostInfo) (errs []error) {
	svc.mu.Lock()
	defer svc.mu.Unlock()

	if svc.state > statePrepare {
		errs = append(errs, alreadyOpened)
		return
	}

	hostNames, vhostParams := info.parse()

	errs = svc.params.validate(vhostParams)
	if len(errs) > 0 {
		return
	}
	svc.params = append(svc.params, vhostParams...)

	vhost := newVhost(info.Cert, hostNames, info.Handler)
	svc.vhosts = append(svc.vhosts, vhost)

	svc.addVhostToServers(vhost, vhostParams)

	return
}

func (svc *Service) openListeners() (errs []error) {
	for _, listener := range svc.listeners {
		err := listener.open()
		if err != nil {
			errs = append(errs, err)
		}
	}

	return
}

func (svc *Service) openServers() (errs []error) {
	chServeErr := make(chan error)

	go func() {
		wg := &sync.WaitGroup{}
		for _, listener := range svc.listeners {
			wg.Add(1)
			l := listener
			go func() {
				err := l.server.open(l)
				if err != nil {
					chServeErr <- err
				}
				wg.Done()
			}()
		}
		wg.Wait()
		close(chServeErr)
	}()

	for err := range chServeErr {
		errs = append(errs, err)
	}

	return
}

func (svc *Service) Open() (errs []error) {
	svc.mu.Lock()
	if svc.state >= stateOpened {
		svc.mu.Unlock()
		errs = append(errs, alreadyOpened)
		return
	}
	svc.state = stateOpened
	svc.mu.Unlock()

	for _, s := range svc.servers {
		s.updateDefaultVhost()
		s.updateHttpServerTLSConfig()
		s.updateHttpServerHandler()
	}

	defer svc.Close()

	errs = svc.openListeners()
	if len(errs) > 0 {
		return
	}

	errs = svc.openServers()
	return
}

func (svc *Service) Close() {
	svc.mu.Lock()
	if svc.state >= stateClosed {
		svc.mu.Unlock()
		return
	}
	svc.state = stateClosed
	svc.mu.Unlock()

	for _, listener := range svc.listeners {
		if listener.server != nil {
			listener.server.close()
		}
		listener.close()
	}
}

func (svc *Service) GetAccessibleURLs(includeLoopback bool) [][]string {
	gotIPList := false
	var ipv46s, ipv4s, ipv6s []string

	vhUrls := make(map[*vhost][]string)

	for _, listener := range svc.listeners {
		server := listener.server
		defaultVh := server.getDefaultVhost()

		port := ""
		if !isDefaultPort(listener.port, server.useTLS) {
			port = listener.port
		}

		for _, vh := range server.vhosts {
			if listener.proto == unix {
				url := "unix:" + listener.port
				vhUrls[vh] = append(vhUrls[vh], url)
				continue
			}
			for _, hostname := range vh.hostNames {
				var url string
				if server.useTLS {
					url = httpsUrl
				} else {
					url = httpUrl
				}
				url = url + hostname + port
				vhUrls[vh] = append(vhUrls[vh], url)
			}
			if vh == defaultVh {
				var url string
				if server.useTLS {
					url = httpsUrl
				} else {
					url = httpUrl
				}
				if len(listener.ip) > 0 {
					url = url + listener.ip + port
					vhUrls[vh] = append(vhUrls[vh], url)
				} else {
					if !gotIPList {
						gotIPList = true
						ipv46s, ipv4s, ipv6s = getAllIfaceIPs(includeLoopback)
					}
					var ips []string
					switch listener.proto {
					case tcp46:
						ips = ipv46s
					case tcp4:
						ips = ipv4s
					case tcp6:
						ips = ipv6s
					}
					for _, ip := range ips {
						ipUrl := url + ip + port
						vhUrls[vh] = append(vhUrls[vh], ipUrl)
					}
				}
			}
		}
	}

	vhostsUrls := make([][]string, len(svc.vhosts))
	for i, vhost := range svc.vhosts {
		vhostsUrls[i] = vhUrls[vhost]
	}
	return vhostsUrls
}
