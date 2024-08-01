package goVirtualHost

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
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
	loadedCerts, errs := LoadCertificatesFromPairs(vh.certKeyPaths)

	vh.loadedCerts = vh.loadedCerts[0:0]
	vh.loadedCerts = append(vh.loadedCerts, loadedCerts...)
	vh.loadedCerts = append(vh.loadedCerts, vh.certs...)

	return errs
}

func (vh *vhost) lookupCertificate(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	certLen := len(vh.loadedCerts)
	if certLen == 1 {
		return vh.loadedCerts[0], nil
	}

	for _, cert := range vh.loadedCerts {
		err := hello.SupportsCertificate(cert)
		if err == nil {
			return cert, err
		}
	}

	for _, cert := range vh.loadedCerts {
		if cert.Leaf == nil {
			cert.Leaf, _ = x509.ParseCertificate(cert.Certificate[0])
			if cert.Leaf == nil {
				continue
			}
		}
		err := cert.Leaf.VerifyHostname(hello.ServerName)
		if err == nil {
			return cert, err
		}
	}

	if certLen > 0 {
		return vh.loadedCerts[0], nil
	}

	return nil, errors.New("cannot find proper certificate for " + hello.ServerName)
}
