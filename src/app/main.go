package app

import (
	"../goVirtualHost"
	"../param"
	"../serverErrHandler"
	"../serverLog"
	"../vhostMux"
	"crypto/tls"
	"os"
)

type App struct {
	vhostSvc   *goVirtualHost.Service
	vhostMuxes []*vhostMux.VhostMux
}

func (app *App) Open() {
	errors := app.vhostSvc.Open()
	for _, err := range errors {
		serverErrHandler.CheckError(err)
	}
}

func (app *App) Close() {
	for _, vhMux := range app.vhostMuxes {
		vhMux.Close()
	}

	app.vhostSvc.Close()
}

func (app *App) ReOpenLog() {
	for _, vhMux := range app.vhostMuxes {
		vhMux.ReOpenLog()
	}
}

func NewApp(params []*param.Param) *App {
	vhSvc := goVirtualHost.NewService()
	vhMuxes := make([]*vhostMux.VhostMux, 0, len(params))

	for _, p := range params {
		// logger
		logger := serverLog.NewLogger(p.AccessLog, p.ErrorLog)
		errors := logger.Open()
		serverErrHandler.CheckFatal(errors...)

		// ErrHandler
		errHandler := serverErrHandler.NewErrHandler(logger)

		// ServeMux
		vhMux := vhostMux.NewServeMux(p, logger, errHandler)
		vhMuxes = append(vhMuxes, vhMux)

		// cert
		var cert *tls.Certificate
		if len(p.Cert) > 0 && len(p.Key) > 0 {
			c, err := tls.LoadX509KeyPair(p.Cert, p.Key)
			if err != nil {
				serverErrHandler.CheckFatal(err)
				logger.LogErrors(err)
			} else {
				cert = &c
			}
		}

		// init vhost
		errors = vhSvc.Add(&goVirtualHost.HostInfo{
			Listens:      p.Listens,
			ListensPlain: p.ListensPlain,
			ListensTLS:   p.ListensTLS,
			Cert:         cert,
			HostNames:    p.HostNames,
			Handler:      vhMux.ServeMux,
		})
		if len(errors) > 0 {
			serverErrHandler.CheckFatal(errors...)
			logger.LogErrors(errors...)
			os.Exit(1)
		}
	}

	return &App{
		vhostSvc:   vhSvc,
		vhostMuxes: vhMuxes,
	}
}
