package app

import (
	"../vhost"
	"crypto/tls"
	"net"
	"net/http"
)

type ListenItem struct {
	proto        string
	addr         string
	useTLS       bool
	hostnames    []string
	certs        []tls.Certificate
	handler      http.Handler
	listener     net.Listener
	server       *http.Server
	vhosts       []*vhost.VHost
	defaultVHost *vhost.VHost
}

type Listens []*ListenItem

func (l *ListenItem) containsHostname(hostname string) bool {
	for _, n := range l.hostnames {
		if n == hostname {
			return true
		}
	}

	return false
}

func (ls Listens) findItemByAddr(addr string) *ListenItem {
	for _, l := range ls {
		if l.addr == addr {
			return l
		}
	}

	return nil
}

func (ls Listens) findItemByAddrHostname(addr, hostname string) *ListenItem {
	for _, l := range ls {
		if l.addr == addr && l.containsHostname(hostname) {
			return l
		}
	}

	return nil
}
