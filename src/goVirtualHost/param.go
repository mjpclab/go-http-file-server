package goVirtualHost

import (
	"errors"
	"fmt"
)

var CertificateNotFound = errors.New("certificate not found for TLS listens")

func (param *param) hasHostNames(checkHostNames []string) bool {
	if len(param.hostNames) == 0 || len(checkHostNames) == 0 {
		return false
	}

	for _, checkHostName := range checkHostNames {
		if contains(param.hostNames, checkHostName) {
			return true
		}
	}

	return false
}

func (param *param) hasHostName(checkHostName string) bool {
	if len(param.hostNames) == 0 || len(checkHostName) == 0 {
		return false
	}

	if contains(param.hostNames, checkHostName) {
		return true
	}

	return false
}

func (param *param) stacksEqual(other *param) bool {
	return param.proto == other.proto && param.ip == other.ip && param.port == other.port
}

func (param *param) validate() (errs []error) {
	if param.useTLS && len(param.certs) == 0 {
		err := wrapError(CertificateNotFound, fmt.Sprintf("certificate not found for TLS listens: %+v", param))
		errs = append(errs, err)
	}

	return
}
