package goVirtualHost

import "fmt"

func (params params) validateParam(param *param) (errs []error) {
	proto := param.proto
	addr := param.addr
	hostnames := param.hostNames

	for _, ownParam := range params {
		if ownParam == param {
			continue
		}

		if ownParam.proto != proto || ownParam.addr != addr {
			continue
		}

		ownUseTLS := ownParam.cert != nil
		inputUseTLS := param.cert != nil
		if ownUseTLS != inputUseTLS {
			err := fmt.Errorf("cannot serve for both Plain and TLS mode: %+v", param)
			errs = append(errs, err)
		}

		if (len(hostnames) == 0 && len(ownParam.hostNames) == 0) || (ownParam.hasHostNames(hostnames)) {
			err := fmt.Errorf("duplicated address and hostname: %+v", param)
			errs = append(errs, err)
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

		es = params.validateParam(p)
		if len(es) > 0 {
			errs = append(errs, es...)
		}

		es = inputs.validateParam(p)
		if len(es) > 0 {
			errs = append(errs, es...)
		}
	}

	return
}
