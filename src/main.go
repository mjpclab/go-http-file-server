package main

import (
	"./param"
	"./server"
	"./serverErrorHandler"
	"./serverLog"
	"os"
	"os/signal"
)

var p *param.Param
var logger *serverLog.Logger

func init() {
	p = param.Parse()

	var err error
	logger, err = serverLog.NewLogger(p.AccessLog, p.ErrorLog)
	if !serverErrorHandler.CheckFatal(err) {
		serverErrorHandler.SetLogger(logger)
	}
}

func cleanup() {
	logger.Close()
}

func cleanupOnInterrupt() {
	// trap SIGINT
	chSignal := make(chan os.Signal)
	signal.Notify(chSignal, os.Interrupt)
	go func() {
		<-chSignal
		cleanup()
		os.Exit(0)
	}()
}

func main() {
	defer cleanup()

	cleanupOnInterrupt()

	s := server.NewServer(p, logger)
	s.ListenAndServe()
}
