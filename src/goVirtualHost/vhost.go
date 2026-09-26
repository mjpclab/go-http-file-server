package goVirtualHost

import (
	"net/http"
	"strings"
)

func newVhost(hostNames []string, certKeyPaths certKeyPairs, vhCerts certs, handler http.Handler) *vhost {
	loadedCerts := make(certs, 0, len(certKeyPaths)+len(vhCerts))
	loadedCerts = append(loadedCerts, vhCerts...)

	vhost := &vhost{
		hostNames:    hostNames,
		certKeyPaths: certKeyPaths,
		certs:        vhCerts,
		loadedCerts:  loadedCerts,
		handler:      handler,
	}

	return vhost
}

func (vh *vhost) matchHostName(name string) bool {
	reqHostName := strings.ToLower(name)
	for _, hostname := range vh.hostNames {
		if hostname == reqHostName {
			return true
		}
		if len(hostname) > 1 {
			if hostname[0] == '.' && strings.HasSuffix(reqHostName, hostname) {
				return true
			} else if hostname[len(hostname)-1] == '.' && strings.HasPrefix(reqHostName, hostname) {
				return true
			}
		}
	}
	return false
}

func (vh *vhost) loadCertificates() []error {
	fileCerts, errs := LoadCertificatesFromPairs(vh.certKeyPaths)
	loadedCerts := certs(fileCerts)
	loadedCerts = append(loadedCerts, vh.certs...)

	vh.loadedCerts = loadedCerts

	return errs
}
