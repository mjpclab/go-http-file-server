package goVirtualHost

import (
	"context"
	"errors"
	"sync"
)

var alreadyOpened = errors.New("already opened")

func NewService() *Service {
	service := &Service{
		state:       statePrepare,
		listenables: listenables{},
		serveables:  serveables{},
		vhosts:      vhosts{},
	}

	return service
}

func (svc *Service) addVhostToServeable(vhost *vhost, params params) {
	for _, param := range params {
		var l *listenable
		var s *serveable

		l = svc.listenables.find(param.proto, param.ip, param.port)
		if l != nil {
			s = l.serveable
		} else {
			s = newServeable(param.useTLS)
			l = newListenable(param.proto, param.ip, param.port)
			l.serveable = s

			svc.listenables = append(svc.listenables, l)
			svc.serveables = append(svc.serveables, s)
		}

		// serveable -> vhost
		s.vhosts = append(s.vhosts, vhost)
	}
}

func (svc *Service) Add(info *HostInfo) (errs, warns []error) {
	svc.mu.Lock()
	defer svc.mu.Unlock()

	if svc.state > statePrepare {
		errs = append(errs, alreadyOpened)
		return
	}

	vhostParams, hostNames, certKeyPaths, certs := info.parse()

	errs, warns = svc.params.validate(vhostParams)
	if len(errs) > 0 {
		return
	}
	svc.params = append(svc.params, vhostParams...)

	vhost := newVhost(hostNames, certKeyPaths, certs, info.Handler)
	svc.vhosts = append(svc.vhosts, vhost)

	svc.addVhostToServeable(vhost, vhostParams)

	return
}

func (svc *Service) openListeners() (errs []error) {
	for _, l := range svc.listenables {
		err := l.open()
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
		for _, lis := range svc.listenables {
			wg.Add(1)
			l := lis
			go func() {
				err := l.serveable.open(l)
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

	svc.params = nil // release unused data

	for _, s := range svc.serveables {
		es := s.init()
		errs = append(errs, es...)
	}
	if len(errs) > 0 {
		return
	}

	defer svc.Close()

	errs = svc.openListeners()
	if len(errs) > 0 {
		return
	}

	errs = svc.openServers()
	return
}

func (svc *Service) ReloadCertificates() (errs []error) {
	for _, s := range svc.serveables {
		es := s.loadCertificates()
		errs = append(errs, es...)
	}
	return
}

func (svc *Service) Shutdown(ctx context.Context) {
	svc.mu.Lock()
	if svc.state >= stateClosed {
		svc.mu.Unlock()
		return
	}
	svc.state = stateClosed
	svc.mu.Unlock()

	wg := &sync.WaitGroup{}
	for _, lis := range svc.listenables {
		l := lis

		wg.Add(1)
		go func() {
			if l.serveable != nil {
				l.serveable.shutdown(ctx)
			}
			l.close()
			wg.Done()
		}()
	}
	wg.Wait()
}

func (svc *Service) Close() {
	svc.mu.Lock()
	if svc.state >= stateClosed {
		svc.mu.Unlock()
		return
	}
	svc.state = stateClosed
	svc.mu.Unlock()

	for _, l := range svc.listenables {
		if l.serveable != nil {
			l.serveable.close()
		}
		l.close()
	}
}

func (svc *Service) GetAccessibleURLs(includeLoopback bool) [][]string {
	gotIPList := false
	var ipv46s, ipv4s, ipv6s []string
	var allHostNameUrls []string
	vhUrls := make(map[*vhost][]string)

	for _, l := range svc.listenables {
		s := l.serveable
		s.updateDefaultVhost()

		port := ""
		if !isDefaultPort(l.port, s.useTLS) {
			port = l.port
		}

		for _, vh := range s.vhosts {
			if l.proto == unix {
				url := "unix:" + l.port
				vhUrls[vh] = append(vhUrls[vh], url)
				continue
			}
			for _, hostname := range vh.hostNames {
				var url string
				if s.useTLS {
					url = httpsUrl
				} else {
					url = httpUrl
				}
				url = url + hostname + port
				if !contains(allHostNameUrls, url) {
					allHostNameUrls = append(allHostNameUrls, url)
					vhUrls[vh] = append(vhUrls[vh], url)
				}
			}

			if vh != s.defaultVhost {
				continue
			}
			var url string
			if s.useTLS {
				url = httpsUrl
			} else {
				url = httpUrl
			}
			if len(l.ip) > 0 {
				url = url + l.ip + port
				vhUrls[vh] = append(vhUrls[vh], url)
				continue
			}

			if !gotIPList {
				gotIPList = true
				ipv46s, ipv4s, ipv6s = getAllIfaceIPs(includeLoopback)
			}
			var ips []string
			switch l.proto {
			case tcp46:
				ips = ipv46s
			case tcp4:
				ips = ipv4s
			case tcp6:
				ips = ipv6s
			}
			for _, ip := range ips {
				if ipVh := s.lookupVhost(ip); ipVh == vh {
					ipUrl := url + ip + port
					vhUrls[vh] = append(vhUrls[vh], ipUrl)
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
