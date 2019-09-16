package main

import (
	"./param"
	"./vhost"
	"os"
	"os/signal"
	"syscall"
)

var h *vhost.VHost

func cleanupOnInterrupt() {
	chSignal := make(chan os.Signal)
	signal.Notify(chSignal, syscall.SIGINT)

	go func() {
		<-chSignal
		h.Close()
		os.Exit(0)
	}()
}

func reOpenLogOnHup() {
	chSignal := make(chan os.Signal)
	signal.Notify(chSignal, syscall.SIGHUP)

	go func() {
		for range chSignal {
			h.ReOpenLog()
		}
	}()
}

func main() {
	cleanupOnInterrupt()

	p := param.ParseCli()
	h = vhost.NewVHost(p)
	if h != nil {
		reOpenLogOnHup()
		h.Open()
		defer h.Close()
	}
}
