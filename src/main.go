package main

import (
	"./param"
	"./vhost"
	"os"
	"os/signal"
)

var h *vhost.VHost

func cleanupOnInterrupt() {
	// trap SIGINT
	chSignal := make(chan os.Signal)
	signal.Notify(chSignal, os.Interrupt)
	go func() {
		<-chSignal
		h.Close()
		os.Exit(0)
	}()
}

func main() {
	cleanupOnInterrupt()

	p := param.ParseCli()
	h = vhost.NewVHost(p)
	if h != nil {
		h.Open()
		defer h.Close()
	}
}
