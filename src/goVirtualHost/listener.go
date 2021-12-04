package goVirtualHost

import (
	"net"
	"os"
)

func newListener(proto, ip, port string) *listener {
	listener := &listener{
		proto: proto,
		ip:    ip,
		port:  port,
	}

	return listener
}

func (listener *listener) open() error {
	addr := listener.ip + listener.port
	if listener.proto == "unix" {
		sockInfo, _ := os.Lstat(addr)
		if sockInfo != nil && (sockInfo.Mode()&os.ModeSocket != 0) {
			os.Remove(addr)
		}
	}

	netListener, err := net.Listen(listener.proto, addr)
	listener.netListener = netListener

	if listener.proto == "unix" && err == nil {
		os.Chmod(addr, 0777)
	}

	return err
}

func (listener *listener) close() error {
	if listener.netListener == nil {
		return nil
	}

	err := listener.netListener.Close()
	listener.netListener = nil
	return err
}
