package main

import (
	"./param"
	"./server"
	"os"
	"os/signal"
)

var s *server.Server

func cleanupOnInterrupt() {
	// trap SIGINT
	chSignal := make(chan os.Signal)
	signal.Notify(chSignal, os.Interrupt)
	go func() {
		<-chSignal
		if s != nil {
			s.Close()
		}
		os.Exit(0)
	}()
}

func main() {
	cleanupOnInterrupt()

	p := param.ParseCli()
	s = server.NewServer(p)
	defer s.Close()
	s.ListenAndServe()
}
