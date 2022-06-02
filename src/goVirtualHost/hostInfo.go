package goVirtualHost

import "crypto/tls"

func (info *HostInfo) toParam(listen string, useTLS bool) *param {
	proto, ip, port := splitListen(listen, false)
	var certs []tls.Certificate
	if useTLS {
		certs = info.Certs
	}

	param := &param{
		proto:  proto,
		ip:     ip,
		port:   port,
		useTLS: useTLS,
		certs:  certs,
	}

	return param
}

func (info *HostInfo) parse() (hostNames []string, params params, certs certs) {
	hostNames = normalizeHostNames(info.HostNames)

	useTLSForListen := len(info.Certs) > 0
	for _, listen := range info.Listens {
		param := info.toParam(listen, useTLSForListen)
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

	if (useTLSForListen && len(info.Listens) > 0) || len(info.ListensTLS) > 0 {
		certs = append(info.Certs)
	}

	return
}
