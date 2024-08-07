package goVirtualHost

import (
	"crypto/tls"
	"errors"
)

var MissingCertFileAndKeyFile = errors.New("missing certificate file and key file")
var MissingCertFile = errors.New("missing certificate file")
var MissingKeyFile = errors.New("missing key file")
var CertKeyFileCountNotMatch = errors.New("certificate file count and key file count not match")

func LoadCertificate(certFile, keyFile string) (cert *tls.Certificate, err error) {
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

	var tlsCert tls.Certificate
	tlsCert, err = tls.LoadX509KeyPair(certFile, keyFile)
	cert = &tlsCert
	return
}

func LoadCertificates(certFiles, keyFiles []string) (certs []*tls.Certificate, errs []error) {
	certLen := len(certFiles)
	if certLen != len(keyFiles) {
		errs = append(errs, CertKeyFileCountNotMatch)
		return
	}
	if certLen == 0 {
		return
	}

	certs = make([]*tls.Certificate, 0, certLen)
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

func LoadCertificatesFromPairs(certKeyFilePairs [][2]string) (certs []*tls.Certificate, errs []error) {
	certLen := len(certKeyFilePairs)
	if certLen == 0 {
		return
	}

	certs = make([]*tls.Certificate, 0, certLen)
	for i := 0; i < certLen; i++ {
		cert, err := LoadCertificate(certKeyFilePairs[i][0], certKeyFilePairs[i][1])
		if err != nil {
			errs = append(errs, err)
		} else {
			certs = append(certs, cert)
		}
	}
	return
}

func CertsKeysToPairs(certFiles, keyFiles []string) (certKeyFilePairs [][2]string, errs []error) {
	certLen := len(certFiles)
	if certLen != len(keyFiles) {
		errs = []error{CertKeyFileCountNotMatch}
		return
	}
	if certLen == 0 {
		return
	}

	certKeyFilePairs = make([][2]string, 0, certLen)
	for i := 0; i < certLen; i++ {
		certKeyFilePairs = append(certKeyFilePairs, [2]string{certFiles[i], keyFiles[i]})
	}
	return
}
