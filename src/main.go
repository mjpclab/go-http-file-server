package main

import (
	"./param"
	"./server"
	"./serverError"
	"./serverLog"
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

	s := server.NewServer(p, logger)
	s.ListenAndServe()
}
