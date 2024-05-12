package goVirtualHost

import (
	"net"
	"os"
)

func newListenable(proto, ip, port string) *listenable {
	return &listenable{
		proto: proto,
		ip:    ip,
		port:  port,
	}
}

func (l *listenable) open() error {
	addr := l.ip + l.port
	if l.proto == "unix" {
		sockInfo, _ := os.Lstat(addr)
		if sockInfo != nil && (sockInfo.Mode()&os.ModeSocket != 0) {
			os.Remove(addr)
		}
	}

	listener, err := net.Listen(l.proto, addr)
	l.listener = listener

	if l.proto == "unix" && err == nil {
		os.Chmod(addr, 0660)
	}

	return err
}

func (l *listenable) close() error {
	if l.listener == nil {
		return nil
	}

	err := l.listener.Close()
	l.listener = nil
	return err
}
