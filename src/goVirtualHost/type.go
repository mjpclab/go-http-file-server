package goVirtualHost

import (
	"crypto/tls"
	"net"
	"net/http"
	"sync"
)

// init host info
type HostInfo struct {
	Listens      []string
	ListensPlain []string
	ListensTLS   []string
	Cert         *tls.Certificate
	HostNames    []string
	Handler      http.Handler
}

// normalized HostInfo Param
type param struct {
	proto     string
	addr      string
	useTLS    bool
	cert      *tls.Certificate
	hostNames []string
	handler   http.Handler
}

type params []*param

// wrapper of net.Listener
type listener struct {
	proto       string
	addr        string
	netListener net.Listener
	server      *server
}

type listeners []*listener

// wrapper for http.Server
type server struct {
	useTLS       bool
	vhosts       vhosts
	defaultVhost *vhost
	httpServer   *http.Server
}

type servers []*server

// virtual host
type vhost struct {
	cert      *tls.Certificate
	hostNames []string
	handler   http.Handler
}

type vhosts []*vhost

// service

type state int

const (
	statePrepare state = iota
	stateOpened
	stateClosed
)

type Service struct {
	mu        sync.Mutex
	state     state
	params    params
	listeners listeners
	servers   servers
	vhosts    vhosts
}
