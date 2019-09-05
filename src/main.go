package main

import (
	"./param"
	"./server"
	"./serverError"
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
	if !serverError.CheckFatal(err) {
		serverError.SetLogger(logger)
	}
}

func main() {
	defer logger.Close()

	// trap SIGINT
	chInter := make(chan os.Signal)
	signal.Notify(chInter, os.Interrupt)
	go func() {
		<-chInter
		logger.Close()
		os.Exit(0)
	}()

	// start server
	s := server.NewServer(p, logger)
	s.ListenAndServe()
}
