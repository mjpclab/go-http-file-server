package app

import (
	"../goVirtualHost"
	"../param"
	"../serverErrHandler"
	"../serverHandler"
	"../serverLog"
	"../tpl"
	"../util"
	"../vhost"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

type App struct {
	vhostSvc   *goVirtualHost.Service
	logFileMan *serverLog.FileMan
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
	writePidFile()

	verbose := !util.GetBoolEnv("GHFS_QUIET")

	if serverHandler.TryEnableWSL1Fix() && verbose {
		ttyFile, teardownTtyFile := util.GetTTYFile()
		fmt.Fprintln(ttyFile, "WSL 1 compatible mode enabled")
		teardownTtyFile()
	}

	vhSvc := goVirtualHost.NewService()
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

		// vHost Handler
		vhHandler := vhost.NewHandler(p, logger, errHandler, theme)

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
			Handler:      vhHandler,
		})
		if len(errors) > 0 {
			serverErrHandler.CheckFatal(errors...)
			logger.LogErrors(errors...)
			os.Exit(1)
		}
	}

	if verbose {
		go printAccessibleURLs(vhSvc)
	}

	return &App{
		vhostSvc:   vhSvc,
		logFileMan: logFileMan,
	}
}

func writePidFile() {
	pidFilename := os.Getenv("GHFS_PID_FILE")
	if len(pidFilename) == 0 {
		return
	}

	pidContent := strconv.Itoa(os.Getpid())
	pidFile, err := os.OpenFile(pidFilename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if !serverErrHandler.CheckFatal(err) {
		_, err := pidFile.WriteString(pidContent)
		err2 := pidFile.Close()
		serverErrHandler.CheckFatal(err, err2)
	}
}
