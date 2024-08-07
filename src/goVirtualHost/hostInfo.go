package goVirtualHost

func (info *HostInfo) toParam(listen string, useTLS bool) *param {
	proto, ip, port := splitListen(listen, false)
	var certKeyPaths certKeyPairs
	var certs certs
	if useTLS {
		certKeyPaths = info.CertKeyPaths
		certs = info.Certs
	}

	param := &param{
		proto:        proto,
		ip:           ip,
		port:         port,
		useTLS:       useTLS,
		certKeyPaths: certKeyPaths,
		certs:        certs,
	}

	return param
}

func (info *HostInfo) parse() (params params, hostNames []string, certKeyPaths certKeyPairs, certs certs) {
	hostNames = normalizeHostNames(info.HostNames)

	useTLSForListen := len(info.CertKeyPaths)+len(info.Certs) > 0
	if (useTLSForListen && len(info.Listens) > 0) || len(info.ListensTLS) > 0 {
		certKeyPaths = info.CertKeyPaths
		certs = info.Certs
	}

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

	return
}
