package app

import (
	"context"
	"errors"
	"io"
	"net/http"
	"time"

	"mjpclab.dev/ghfs/src/goVirtualHost"
	"mjpclab.dev/ghfs/src/param"
	"mjpclab.dev/ghfs/src/serverHandler"
	"mjpclab.dev/ghfs/src/serverLog"
	"mjpclab.dev/ghfs/src/tpl/theme"
)

type App struct {
	params   param.Params
	vhostSvc *goVirtualHost.Service
	logMan   serverLog.Man
}

type fnGetLogger func(paramIndex int, param *param.Param) (logger *serverLog.Logger, errs []error)

var errWriterNumberNotMatchParams = errors.New("number of writers not equal to params")

func (app *App) Open() []error {
	errs := app.vhostSvc.Open()
	es := make([]error, 0, len(errs))
	for i := range errs {
		if errs[i] != http.ErrServerClosed {
			es = append(es, errs[i])
		}
	}
	return es
}

func (app *App) Close() {
	app.vhostSvc.Close()
	app.logMan.Close()
}

func (app *App) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*300)
	app.vhostSvc.Shutdown(ctx)
	cancel()

	app.logMan.Close()
}

func (app *App) ReOpen() []error {
	return app.logMan.ReOpen()
}

func (app *App) ReLoadCertificates() []error {
	return app.vhostSvc.ReloadCertificates()
}

func newLogApp(params param.Params, logMan serverLog.Man, fnGetLogger fnGetLogger) (*App, []error) {
	vhSvc := goVirtualHost.NewService()
	themePool := make(map[string]theme.Theme)

	for i, p := range params {
		// logger
		logger, errs := fnGetLogger(i, p)
		if len(errs) > 0 {
			return nil, errs
		}

		// theme
		var themeInst theme.Theme
		if len(p.ThemeDir) > 0 {
			themeInst = theme.DirTheme(p.ThemeDir)
		} else if len(p.Theme) > 0 {
			themeInst, errs = loadTheme(p.Theme, themePool)
			if len(errs) > 0 {
				logger.LogErrors(errs...)
				return nil, errs
			}
		}

		// vHost Handler
		vhHandler, errs := serverHandler.NewVhostHandler(p, logger, themeInst)
		if len(errs) > 0 {
			logger.LogErrors(errs...)
			return nil, errs
		}

		// init vhost
		listens := p.Listens
		if len(listens) == 0 && len(p.ListensPlain) == 0 && len(p.ListensTLS) == 0 {
			if len(p.CertKeyPaths) == 0 {
				listens = []string{":80"}
			} else {
				listens = []string{":443"}
			}
		}

		var warns []error
		errs, warns = vhSvc.Add(&goVirtualHost.HostInfo{
			Listens:      listens,
			ListensPlain: p.ListensPlain,
			ListensTLS:   p.ListensTLS,
			CertKeyPaths: p.CertKeyPaths,
			HostNames:    p.HostNames,
			Handler:      vhHandler,
		})
		if len(warns) > 0 {
			logger.LogErrors(warns...)
		}
		if len(errs) > 0 {
			logger.LogErrors(errs...)
			return nil, errs
		}
	}

	return &App{
		params:   params,
		vhostSvc: vhSvc,
		logMan:   logMan,
	}, nil
}

func NewApp(params param.Params) (*App, []error) {
	logMan := serverLog.NewHybridMan()
	app, errs := newLogApp(params, logMan, func(paramIndex int, param *param.Param) (*serverLog.Logger, []error) {
		return logMan.NewLogger(param.AccessLog, param.ErrorLog)
	})
	if len(errs) > 0 {
		logMan.Close()
		return nil, errs
	}
	return app, errs
}

func NewWriterLogApp(params param.Params, writers [][2]io.Writer) (*App, []error) {
	if len(writers) != len(params) {
		return nil, []error{errWriterNumberNotMatchParams}
	}

	logMan := serverLog.NewWriterMan()
	app, errs := newLogApp(params, logMan, func(paramIndex int, param *param.Param) (*serverLog.Logger, []error) {
		logger := logMan.NewLogger(writers[paramIndex][0], writers[paramIndex][1])
		return logger, nil
	})
	if len(errs) > 0 {
		logMan.Close()
		return nil, errs
	}
	return app, errs
}
