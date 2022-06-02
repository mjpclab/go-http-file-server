package goVirtualHost

import (
	"crypto/tls"
	"errors"
)

var MissingCertFileAndKeyFile = errors.New("missing certificate file and key file")
var MissingCertFile = errors.New("missing certificate file")
var MissingKeyFile = errors.New("missing key file")
var CertKeyFileCountNotMatch = errors.New("certificate file count and key file count not match")

func LoadCertificate(certFile, keyFile string) (cert tls.Certificate, err error) {
	if len(certFile) == 0 && len(keyFile) == 0 {
		err = MissingCertFileAndKeyFile
		return
	} else if len(certFile) == 0 {
		err = MissingCertFile
		return
	} else if len(keyFile) == 0 {
		err = MissingKeyFile
		return
	}

	cert, err = tls.LoadX509KeyPair(certFile, keyFile)
	return
}

func LoadCertificates(certFiles, keyFiles []string) (certs []tls.Certificate, errs []error) {
	certLen := len(certFiles)
	if certLen != len(keyFiles) {
		errs = append(errs, CertKeyFileCountNotMatch)
		return
	}

	if certLen == 0 {
		return
	}

	certs = make([]tls.Certificate, 0, certLen)
	for i := 0; i < certLen; i++ {
		cert, err := LoadCertificate(certFiles[i], keyFiles[i])
		if err != nil {
			errs = append(errs, err)
		} else {
			certs = append(certs, cert)
		}
	}

	return
}

func LoadCertificatesFromEntries(certKeyFileEntries [][2]string) (certs []tls.Certificate, errs []error) {
	certLen := len(certKeyFileEntries)
	if certLen == 0 {
		return
	}

	certs = make([]tls.Certificate, 0, certLen)
	for i := 0; i < certLen; i++ {
		cert, err := LoadCertificate(certKeyFileEntries[i][0], certKeyFileEntries[i][1])
		if err != nil {
			errs = append(errs, err)
		} else {
			certs = append(certs, cert)
		}
	}
	return
}
