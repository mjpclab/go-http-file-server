package goVirtualHost

import "crypto/x509"

func (cs certs) makeLeaf() {
	for _, cert := range cs {
		if cert.Leaf == nil && len(cert.Certificate) > 0 {
			if leaf, err := x509.ParseCertificate(cert.Certificate[0]); err == nil {
				cert.Leaf = leaf
			}
		}
	}
}
