package goVirtualHost

import "crypto/tls"

func (info *HostInfo) toParam(listen string, useTLS bool) *param {
	proto, ip, port := splitListen(listen, false)
	var cert *tls.Certificate
	if useTLS {
		cert = info.Cert
	}

	param := &param{
		proto:  proto,
		ip:     ip,
		port:   port,
		useTLS: useTLS,
		cert:   cert,
	}

	return param
}

func (info *HostInfo) parse() (hostNames []string, params params) {
	hostNames = normalizeHostNames(info.HostNames)

	for _, listen := range info.Listens {
		param := info.toParam(listen, info.Cert != nil)
		param.hostNames = hostNames
		params = append(params, param)
	}

	for _, listen := range info.ListensPlain {
		param := info.toParam(listen, false)
		param.hostNames = hostNames
		params = append(params, param)
	}

	for _, listen := range info.ListensTLS {
		param := info.toParam(listen, true)
		param.hostNames = hostNames
		params = append(params, param)
	}

	return
}
