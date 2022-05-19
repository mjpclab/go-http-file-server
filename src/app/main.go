package app

import (
	"../goVirtualHost"
	"../param"
	"../serverErrHandler"
	"../serverLog"
	"../tpl"
	"../util"
	"../vhostHandler"
	"os"
	"path/filepath"
)

type App struct {
	vhostSvc      *goVirtualHost.Service
	vhostHandlers []*vhostHandler.VhostHandler
	logFileMan    *serverLog.FileMan
}

func (app *App) Open() {
	errors := app.vhostSvc.Open()
	serverErrHandler.CheckError(errors...)
}

func (app *App) Close() {
	app.vhostSvc.Close()
	app.logFileMan.Close()
}

func (app *App) ReOpenLog() {
	errors := app.logFileMan.Reopen()
	serverErrHandler.CheckFatal(errors...)
}

func NewApp(params []*param.Param) *App {
	vhSvc := goVirtualHost.NewService()
	vhHandlers := make([]*vhostHandler.VhostHandler, 0, len(params))
	logFileMan := serverLog.NewFileMan()
	themes := make(map[string]tpl.Theme)

	for _, p := range params {
		// logger
		logger, errors := logFileMan.NewLogger(p.AccessLog, p.ErrorLog)
		serverErrHandler.CheckFatal(errors...)

		// ErrHandler
		errHandler := serverErrHandler.NewErrHandler(logger)

		// theme
		var theme tpl.Theme
		if len(p.ThemeDir) > 0 {
			theme = tpl.DirTheme(p.ThemeDir)
		} else if len(p.Theme) == 0 {
			theme = tpl.DefaultTheme
		} else {
			themeKey, err := filepath.Abs(p.Theme)
			serverErrHandler.CheckFatal(err)

			var themeExists bool
			theme, themeExists = themes[themeKey]
			if !themeExists {
				theme, err = tpl.LoadMemTheme(p.Theme)
				serverErrHandler.CheckFatal(err)
				themes[themeKey] = theme
			}
		}

		// vHostMux
		vhHandler := vhostHandler.NewHandler(p, logger, errHandler, theme)
		vhHandlers = append(vhHandlers, vhHandler)

		// init vhost
		listens := p.Listens
		if len(listens) == 0 && len(p.ListensPlain) == 0 && len(p.ListensTLS) == 0 {
			if len(p.Certificates) == 0 {
				listens = []string{":80"}
			} else {
				listens = []string{":443"}
			}
		}

		errors = vhSvc.Add(&goVirtualHost.HostInfo{
			Listens:      listens,
			ListensPlain: p.ListensPlain,
			ListensTLS:   p.ListensTLS,
			Certs:        p.Certificates,
			HostNames:    p.HostNames,
			Handler:      vhHandler.Handler,
		})
		if len(errors) > 0 {
			serverErrHandler.CheckFatal(errors...)
			logger.LogErrors(errors...)
			os.Exit(1)
		}
	}

	if !util.GetBoolEnv("GHFS_QUIET") {
		go printAccessibleURLs(vhSvc)
	}

	return &App{
		vhostSvc:      vhSvc,
		vhostHandlers: vhHandlers,
		logFileMan:    logFileMan,
	}
}
