package vhost

import (
	"../param"
	"../serveMux"
	"../serverErrHandler"
	"../serverLog"
	"context"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

type VHost struct {
	logger     *serverLog.Logger
	errHandler *serverErrHandler.ErrHandler
	listeners  []*listenerItem

	ctx         context.Context
	waitServers sync.WaitGroup
}

type listenerItem struct {
	proto    string
	addr     string
	key      string
	cert     string
	useTLS   bool
	listener net.Listener
	server   *http.Server
}

func removeUnixSocket(addr string) {
	sockInfo, _ := os.Lstat(addr)
	if sockInfo != nil && (sockInfo.Mode()&os.ModeSocket != 0) {
		os.Remove(addr)
	}
}

func startListen(l *listenerItem) (listener net.Listener, err error) {
	if l.proto == "unix" {
		removeUnixSocket(l.addr)
	}

	listener, err = net.Listen(l.proto, l.addr)

	if err == nil && l.proto == "unix" {
		os.Chmod(l.addr, 0777)
	}

	return
}

func (v *VHost) startServe(l *listenerItem) {
	var err error
	if l.useTLS {
		err = l.server.ServeTLS(l.listener, l.cert, l.key)
	} else {
		err = l.server.Serve(l.listener)
	}
	v.errHandler.LogError(err)

	v.waitServers.Done()
}

func (v *VHost) Open() {
	for _, l := range v.listeners {
		v.logger.LogAccessString("start to listen on " + l.proto + ": " + l.addr)

		var err error
		l.listener, err = startListen(l)
		if v.errHandler.LogError(err) {
			v.Close()
			return
		}
	}

	for _, l := range v.listeners {
		v.waitServers.Add(1)
		go v.startServe(l)
	}

	v.waitServers.Wait()
	return
}

func (v *VHost) Close() {
	ctxTimeout, _ := context.WithTimeout(v.ctx, time.Second*3)

	for _, l := range v.listeners {
		if l.server == nil {
			continue
		}
		err := l.server.Shutdown(ctxTimeout)
		v.errHandler.LogError(err)
		l.server = nil
	}

	for _, l := range v.listeners {
		if l.listener == nil {
			continue
		}
		l.listener.Close()
		l.listener = nil
	}

	v.logger.Close()
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

	listeners := make([]*listenerItem, 0, listenersCapacity)
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

	// listener.server
	for _, l := range listeners {
		l.server = &http.Server{Handler: mux}
	}

	// vHost
	vHost := &VHost{
		logger:     logger,
		errHandler: errHandler,
		listeners:  listeners,

		ctx: context.Background(),
	}

	return vHost
}

func appendListeners(listeners []*listenerItem, listens []string, key, cert string, useTLS bool) []*listenerItem {
	for _, listen := range listens {
		proto, addr := splitListen(listen, useTLS)
		listeners = append(listeners, &listenerItem{
			proto:  proto,
			addr:   addr,
			key:    key,
			cert:   cert,
			useTLS: useTLS,
		})
	}

	return listeners
}
