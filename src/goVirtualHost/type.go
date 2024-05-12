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
	Certs        []tls.Certificate
	HostNames    []string
	Handler      http.Handler
}

type certs []tls.Certificate

// normalized HostInfo Param
type param struct {
	proto     string // "tcp", "tcp4", "tcp6"
	ip        string
	port      string
	useTLS    bool
	certs     certs
	hostNames []string
}

type params []*param

// wrapper of net.Listener
type listenable struct {
	proto     string // "tcp", "tcp4", "tcp6"
	ip        string
	port      string
	listener  net.Listener
	serveable *serveable
}

type listenables []*listenable

// wrapper for http.Server
type serveable struct {
	useTLS       bool
	vhosts       vhosts
	defaultVhost *vhost
	server       *http.Server
}

type serveables []*serveable

// virtual host
type vhost struct {
	certs     certs
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
	mu          sync.Mutex
	state       state
	params      params
	listenables listenables
	serveables  serveables
	vhosts      vhosts
}

// ip

type ipAddr struct {
	netIP              net.IP
	version            int
	isGlobalUnicast    bool
	isLinkLocalUnicast bool
	isNonPrivate       bool
	isNonLoopback      bool
}
type ipAddrs []*ipAddr
