package goVirtualHost

import (
	"net"
	"os"
)

func newListener(proto, addr string) *listener {
	listener := &listener{
		proto: proto,
		addr:  addr,
	}

	return listener
}

func (listener *listener) open() error {
	if listener.proto == "unix" {
		sockInfo, _ := os.Lstat(listener.addr)
		if sockInfo != nil && (sockInfo.Mode()&os.ModeSocket != 0) {
			os.Remove(listener.addr)
		}
	}

	netListener, err := net.Listen(listener.proto, listener.addr)
	listener.netListener = netListener

	if listener.proto == "unix" && err == nil {
		os.Chmod(listener.addr, 0777)
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
