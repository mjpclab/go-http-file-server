package goVirtualHost

import (
	"errors"
	"fmt"
)

var ConflictIPAddress = errors.New("conflict IP address")
var ConflictTLSMode = errors.New("cannot serve for both Plain and TLS mode")
var DuplicatedAddressHostname = errors.New("duplicated address and hostname")

func (params params) validateParam(param *param) (errs []error) {
	for _, ownParam := range params {
		if ownParam == param {
			continue
		}

		if ownParam.port == param.port && ownParam.proto != unix && param.proto != unix {
			ipConflict := false
			if ownParam.proto == param.proto {
				if ownParam.ip != param.ip && (ownParam.ip == "" || param.ip == "") {
					ipConflict = true
				}
			} else {
				if ownParam.proto == tcp46 || param.proto == tcp46 {
					ipConflict = true
				}
			}
			if ipConflict {
				err := wrapError(ConflictIPAddress, fmt.Sprintf("conflict IP address: %+v, %+v", ownParam, param))
				errs = append(errs, err)
			}
		}

		if ownParam.proto == param.proto && ownParam.ip == param.ip && ownParam.port == param.port {
			ownUseTLS := ownParam.cert != nil
			useTLS := param.cert != nil
			if ownUseTLS != useTLS {
				err := wrapError(ConflictTLSMode, fmt.Sprintf("cannot serve for both Plain and TLS mode: %+v, %+v", ownParam, param))
				errs = append(errs, err)
			}

			if (len(param.hostNames) == 0 && len(ownParam.hostNames) == 0) || (ownParam.hasHostNames(param.hostNames)) {
				err := wrapError(DuplicatedAddressHostname, fmt.Sprintf("duplicated address and hostname: %+v, %+v", ownParam, param))
				errs = append(errs, err)
			}
		}
	}

	return
}

func (params params) validate(inputs params) (errs []error) {
	for _, p := range inputs {
		es := p.validate()
		if len(es) > 0 {
			errs = append(errs, es...)
		}

		es = inputs.validateParam(p)
		if len(es) > 0 {
			errs = append(errs, es...)
		}

		es = params.validateParam(p)
		if len(es) > 0 {
			errs = append(errs, es...)
		}
	}

	return
}
