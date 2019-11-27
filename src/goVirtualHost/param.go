package goVirtualHost

import "fmt"

func (param *param) hasHostNames(checkHostNames []string) bool {
	if len(param.hostNames) == 0 || len(checkHostNames) == 0 {
		return false
	}

	for _, ownHostName := range param.hostNames {
		for _, checkHostName := range checkHostNames {
			if ownHostName == checkHostName {
				return true
			}
		}
	}
	return false
}

func (param *param) validate() (errs []error) {
	if param.useTLS && param.cert == nil {
		err := fmt.Errorf("certificate not found for TLS listens: %+v", param)
		errs = append(errs, err)
	}

	return
}
