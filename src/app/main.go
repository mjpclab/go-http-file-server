package app

import (
	"context"
	"mjpclab.dev/ghfs/src/goVirtualHost"
	"mjpclab.dev/ghfs/src/param"
	"mjpclab.dev/ghfs/src/serverHandler"
	"mjpclab.dev/ghfs/src/serverLog"
	"mjpclab.dev/ghfs/src/setting"
	"mjpclab.dev/ghfs/src/tpl/defaultTheme"
	"mjpclab.dev/ghfs/src/tpl/theme"
	"net/http"
	"time"
)

type App struct {
	vhostSvc   *goVirtualHost.Service
	logFileMan *serverLog.FileMan
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
	app.logFileMan.Close()
}

func (app *App) Shutdown() {
	ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*100)
	app.vhostSvc.Shutdown(ctx)
	app.logFileMan.Close()
}

func (app *App) ReOpenLog() []error {
	return app.logFileMan.Reopen()
}

func (app *App) ReLoadCertificates() []error {
	return app.vhostSvc.ReloadCertificates()
}

func (app *App) GetAccessibleOrigins(includeLoopback bool) [][]string {
	return app.vhostSvc.GetAccessibleURLs(includeLoopback)
}

func NewApp(params param.Params, settings *setting.Setting) (*App, []error) {
	if len(settings.PidFile) > 0 {
		errs := writePidFile(settings.PidFile)
		if len(errs) > 0 {
			return nil, errs
		}
	}

	vhSvc := goVirtualHost.NewService()
	logFileMan := serverLog.NewFileMan()
	themePool := make(map[string]theme.Theme)

	for _, p := range params {
		// logger
		logger, errs := logFileMan.NewLogger(p.AccessLog, p.ErrorLog)
		if len(errs) > 0 {
			return nil, errs
		}

		// theme
		var themeInst theme.Theme
		if len(p.ThemeDir) > 0 {
			themeInst = theme.DirTheme(p.ThemeDir)
		} else if len(p.Theme) == 0 {
			themeInst = defaultTheme.DefaultTheme
		} else {
			themeInst, errs = loadTheme(p.Theme, themePool)
		}
		if len(errs) > 0 {
			logger.LogErrors(errs...)
			return nil, errs
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

	if !settings.Quiet {
		go printAccessibleURLs(vhSvc, params)
	}

	return &App{
		vhostSvc:   vhSvc,
		logFileMan: logFileMan,
	}, nil
}
