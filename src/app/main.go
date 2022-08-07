package app

import (
	"../goVirtualHost"
	"../param"
	"../serverError"
	"../serverHandler"
	"../serverLog"
	"../setting"
	"../tpl"
	"../util"
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

func NewApp(params []*param.Param, setting *setting.Setting) (*App, []error) {
	if len(setting.PidFile) > 0 {
		errs := writePidFile(setting.PidFile)
		if len(errs) > 0 {
			return nil, errs
		}
	}

	if serverHandler.TryEnableWSL1Fix() && !setting.Quiet {
		ttyFile, teardownTtyFile := util.GetTTYFile()
		ttyFile.WriteString("WSL 1 compatible mode enabled\n")
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
			errs = serverError.AppendError(errs, err)
			if err != nil {
				continue
			}

			var themeExists bool
			theme, themeExists = themes[themeKey]
			if !themeExists {
				theme, err = tpl.LoadMemTheme(p.Theme)
				errs = serverError.AppendError(errs, err)
				if err != nil {
					continue
				}
				themes[themeKey] = theme
			}
		}
		if len(errs) > 0 {
			return nil, errs
		}

		// vHost Handler
		vhHandler, errs := serverHandler.NewVhostHandler(p, logger, theme)
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

	if !setting.Quiet {
		go printAccessibleURLs(vhSvc)
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
