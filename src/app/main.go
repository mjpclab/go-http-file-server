package app

import (
	"context"
	"mjpclab.dev/ghfs/src/goVirtualHost"
	"mjpclab.dev/ghfs/src/param"
	"mjpclab.dev/ghfs/src/serverHandler"
	"mjpclab.dev/ghfs/src/serverLog"
	"mjpclab.dev/ghfs/src/tpl/theme"
	"net/http"
	"time"
)

type App struct {
	params   param.Params
	vhostSvc *goVirtualHost.Service
	logMan   *serverLog.Man
}

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
	app.logMan.CloseFiles()
}

func (app *App) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	app.vhostSvc.Shutdown(ctx)
	cancel()

	app.logMan.CloseFiles()
}

func (app *App) ReOpenLog() []error {
	return app.logMan.ReOpenFiles()
}

func (app *App) ReLoadCertificates() []error {
	return app.vhostSvc.ReloadCertificates()
}

func NewApp(params param.Params) (*App, []error) {
	vhSvc := goVirtualHost.NewService()
	logMan := serverLog.NewMan()
	themePool := make(map[string]theme.Theme)

	for _, p := range params {
		// logger
		logger, errs := logMan.NewLogger(p.AccessLog, p.ErrorLog)
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
