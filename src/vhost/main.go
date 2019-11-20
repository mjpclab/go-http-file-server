package vhost

import (
	"../param"
	"../serveMux"
	"../serverErrHandler"
	"../serverLog"
	"net/http"
	"strings"
)

type VHost struct {
	Listens   []*listenItem
	Mux       *http.ServeMux
	Hostnames []string

	logger     *serverLog.Logger
	errHandler *serverErrHandler.ErrHandler
}

type listenItem struct {
	Proto  string
	Addr   string
	Key    string
	Cert   string
	UseTLS bool
}

func (v *VHost) ReOpenLog() {
	errors := v.logger.ReOpen()
	serverErrHandler.CheckError(errors...)
}

func (v *VHost) Close() {
	v.logger.Close()
}

func (v *VHost) MatchHostname(reqHostname string) bool {
	reqHostname = strings.ToLower(reqHostname)
	for _, hostname := range v.Hostnames {
		if hostname == reqHostname {
			return true
		}
		if len(hostname) > 0 && hostname[0] == '.' && strings.HasSuffix(reqHostname, hostname) {
			return true
		}
	}
	return false
}

func NewVHost(p *param.Param) *VHost {
	// logger
	logger := serverLog.NewLogger(p.AccessLog, p.ErrorLog)
	errors := logger.Open()
	serverErrHandler.CheckFatal(errors...)

	// ErrHandler
	errHandler := serverErrHandler.NewErrHandler(logger)

	// ServeMux
	mux := serveMux.NewServeMux(p, logger, errHandler)

	// hostnames
	hostnames := make([]string, 0, len(p.Hostnames))
	for _, hostname := range p.Hostnames {
		hostnames = append(hostnames, strings.ToLower(hostname))
	}

	// determine can use TLS
	key := p.Key
	cert := p.Cert
	canTLS := len(key) > 0 && len(cert) > 0
	if !canTLS {
		key = ""
		cert = ""
	}

	// init listeners
	listenersCapacity := len(p.Listen) + len(p.ListenPlain)
	if canTLS {
		listenersCapacity += len(p.ListenTLS)
	}

	listeners := make([]*listenItem, 0, listenersCapacity)
	listeners = appendListeners(listeners, p.Listen, key, cert, canTLS)
	listeners = appendListeners(listeners, p.ListenPlain, "", "", false)
	if canTLS {
		listeners = appendListeners(listeners, p.ListenTLS, key, cert, true)
	} else if len(p.ListenTLS) > 0 {
		logger.LogErrorString("key or cert not specified for force-TLS port")
		logger.Close()
		return nil
	}

	if len(listeners) == 0 {
		listeners = appendListeners(listeners, []string{""}, key, cert, canTLS)
	}

	// vHost
	vHost := &VHost{
		Listens:   listeners,
		Mux:       mux,
		Hostnames: hostnames,

		logger:     logger,
		errHandler: errHandler,
	}

	return vHost
}

func appendListeners(listeners []*listenItem, listens []string, key, cert string, useTLS bool) []*listenItem {
	for _, listen := range listens {
		proto, addr := splitListen(listen, useTLS)
		listeners = append(listeners, &listenItem{
			Proto:  proto,
			Addr:   addr,
			Key:    key,
			Cert:   cert,
			UseTLS: useTLS,
		})
	}

	return listeners
}
