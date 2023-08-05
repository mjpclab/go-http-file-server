package app

import (
	"context"
	"mjpclab.dev/ghfs/src/goVirtualHost"
	"mjpclab.dev/ghfs/src/param"
	"mjpclab.dev/ghfs/src/serverError"
	"mjpclab.dev/ghfs/src/serverHandler"
	"mjpclab.dev/ghfs/src/serverLog"
	"mjpclab.dev/ghfs/src/setting"
	"mjpclab.dev/ghfs/src/tpl/defaultTheme"
	"mjpclab.dev/ghfs/src/tpl/theme"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type App struct {
	vhostSvc   *goVirtualHost.Service
	logFileMan *serverLog.FileMan
}

func (app *App) Open() []error {
	return app.vhostSvc.Open()
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

func NewApp(params param.Params, setting *setting.Setting) (*App, []error) {
	if len(setting.PidFile) > 0 {
		errs := writePidFile(setting.PidFile)
		if len(errs) > 0 {
			return nil, errs
		}
	}

	vhSvc := goVirtualHost.NewService()
	logFileMan := serverLog.NewFileMan()
	themes := make(map[string]theme.Theme)

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
			themeKey, err := filepath.Abs(p.Theme)
			errs = serverError.AppendError(errs, err)
			if err != nil {
				continue
			}

			var themeExists bool
			themeInst, themeExists = themes[themeKey]
			if !themeExists {
				themeInst, err = theme.LoadMemTheme(p.Theme)
				errs = serverError.AppendError(errs, err)
				if err != nil {
					continue
				}
				themes[themeKey] = themeInst
			}
		}
		if len(errs) > 0 {
			return nil, errs
		}

		// vHost Handler
		vhHandler, errs := serverHandler.NewVhostHandler(p, logger, themeInst)
		if len(errs) > 0 {
			return nil, errs
		}

		// init vhost
		listens := p.Listens
		if len(listens) == 0 && len(p.ListensPlain) == 0 && len(p.ListensTLS) == 0 {
			if len(p.Certificates) == 0 {
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
			Certs:        p.Certificates,
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

	if !setting.Quiet {
		go printAccessibleURLs(vhSvc, params)
	}

	return &App{
		vhostSvc:   vhSvc,
		logFileMan: logFileMan,
	}, nil
}

func writePidFile(pidFilename string) (errs []error) {
	pidContent := strconv.Itoa(os.Getpid())
	pidFile, err := os.OpenFile(pidFilename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return []error{err}
	}

	_, err = pidFile.WriteString(pidContent)
	errs = serverError.AppendError(errs, err)

	err = pidFile.Close()
	errs = serverError.AppendError(errs, err)

	return
}
