package goVirtualHost

import (
	"crypto/tls"
	"errors"
	"net/http"
	"strings"
)

func newVhost(hostNames []string, certKeyPaths certKeyPairs, vhCerts certs, handler http.Handler) *vhost {
	vhCerts.makeLeaf()

	vhost := &vhost{
		hostNames:    hostNames,
		certKeyPaths: certKeyPaths,
		certs:        vhCerts,
		handler:      handler,
	}
	vhost.loadedCerts.Store(vhCerts)

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
	loadedCerts.makeLeaf()

	loadedCerts = append(loadedCerts, vh.certs...)

	vh.loadedCerts.Store(loadedCerts)

	return errs
}

func (vh *vhost) lookupCertificate(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	certs := vh.loadedCerts.Load().(certs)
	if len(certs) == 0 {
		return nil, errors.New("cannot find proper certificate for " + hello.ServerName)
	}
	if len(certs) == 1 {
		return certs[0], nil
	}

	for _, cert := range certs {
		if cert.Leaf == nil {
			continue
		}
		err := cert.Leaf.VerifyHostname(hello.ServerName)
		if err == nil {
			return cert, err
		}
	}

	return certs[0], nil
}
