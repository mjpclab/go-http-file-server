package app

import (
	"../goVirtualHost"
	"../param"
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

func (app *App) Open() []error {
	return app.vhostSvc.Open()
}

func (app *App) Close() {
	app.vhostSvc.Close()
	app.logFileMan.Close()
}

func (app *App) ReOpenLog() []error {
	return app.logFileMan.Reopen()
}

func NewApp(params []*param.Param) (*App, []error) {
	errs := writePidFile()
	if len(errs) > 0 {
		return nil, errs
	}

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
		logger, errs := logFileMan.NewLogger(p.AccessLog, p.ErrorLog)
		if len(errs) > 0 {
			return nil, errs
		}

		// theme
		var theme tpl.Theme
		if len(p.ThemeDir) > 0 {
			theme = tpl.DirTheme(p.ThemeDir)
		} else if len(p.Theme) == 0 {
			theme = tpl.DefaultTheme
		} else {
			themeKey, err := filepath.Abs(p.Theme)
			if err != nil {
				return nil, []error{err}
			}

			var themeExists bool
			theme, themeExists = themes[themeKey]
			if !themeExists {
				theme, err = tpl.LoadMemTheme(p.Theme)
				if err != nil {
					return nil, []error{err}
				}
				themes[themeKey] = theme
			}
		}

		// vHost Handler
		vhHandler := vhost.NewHandler(p, logger, theme)

		// init vhost
		listens := p.Listens
		if len(listens) == 0 && len(p.ListensPlain) == 0 && len(p.ListensTLS) == 0 {
			if len(p.Certificates) == 0 {
				listens = []string{":80"}
			} else {
				listens = []string{":443"}
			}
		}

		errs = vhSvc.Add(&goVirtualHost.HostInfo{
			Listens:      listens,
			ListensPlain: p.ListensPlain,
			ListensTLS:   p.ListensTLS,
			Certs:        p.Certificates,
			HostNames:    p.HostNames,
			Handler:      vhHandler,
		})
		if len(errs) > 0 {
			logger.LogErrors(errs...)
			return nil, errs
		}
	}

	if verbose {
		go printAccessibleURLs(vhSvc)
	}

	return &App{
		vhostSvc:   vhSvc,
		logFileMan: logFileMan,
	}, nil
}

func writePidFile() (errs []error) {
	pidFilename := os.Getenv("GHFS_PID_FILE")
	if len(pidFilename) == 0 {
		return nil
	}

	pidContent := strconv.Itoa(os.Getpid())
	pidFile, err := os.OpenFile(pidFilename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return []error{err}
	}

	_, err = pidFile.WriteString(pidContent)
	err2 := pidFile.Close()

	if err != nil {
		errs = append(errs, err)
	}
	if err2 != nil {
		errs = append(errs, err2)
	}
	return
}
