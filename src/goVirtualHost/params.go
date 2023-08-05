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

		if ownParam.stacksEqual(param) {
			if ownParam.useTLS != param.useTLS {
				err := wrapError(ConflictTLSMode, fmt.Sprintf("cannot serve for both Plain and TLS mode: %+v, %+v", ownParam, param))
				errs = append(errs, err)
			}
		}
	}

	return
}

func (params params) validateShadows(param *param) (errs []error) {
	if len(params) == 0 {
		return nil
	}

	if len(param.hostNames) == 0 {
		shadowed := false
		for _, ownParam := range params {
			if ownParam.stacksEqual(param) && len(ownParam.hostNames) == 0 {
				shadowed = true
				break
			}
		}
		if !shadowed {
			return nil
		}
	} else {
		for _, hostName := range param.hostNames {
			shadowed := false
			for _, ownParam := range params {
				if ownParam.stacksEqual(param) && ownParam.hasHostName(hostName) {
					shadowed = true
					break
				}
			}
			if !shadowed {
				return nil
			}
		}
	}

	err := wrapError(DuplicatedAddressHostname, fmt.Sprintf("duplicated address and hostname: %+v", param))
	errs = append(errs, err)
	return
}

func (params params) validate(inputs params) (errs, warns []error) {
	for _, p := range inputs {
		es := p.validate()
		errs = append(errs, es...)

		es = inputs.validateParam(p)
		errs = append(errs, es...)

		es = params.validateParam(p)
		errs = append(errs, es...)

		ws := params.validateShadows(p)
		warns = append(warns, ws...)
	}

	return
}
