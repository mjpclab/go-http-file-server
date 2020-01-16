package app

import (
	"../goVirtualHost"
	"../param"
	"../serverErrHandler"
	"../serverLog"
	"../vhostMux"
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

		// init vhost
		listens := p.Listens
		if len(listens) == 0 && len(p.ListensPlain) == 0 && len(p.ListensTLS) == 0 {
			if p.Certificate == nil {
				listens = []string{":80"}
			} else {
				listens = []string{":443"}
			}
		}

		errors = vhSvc.Add(&goVirtualHost.HostInfo{
			Listens:      listens,
			ListensPlain: p.ListensPlain,
			ListensTLS:   p.ListensTLS,
			Cert:         p.Certificate,
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
