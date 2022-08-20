package main

import (
	"errors"
	"mjpclab.dev/ghfs/src/app"
	"mjpclab.dev/ghfs/src/param"
	"mjpclab.dev/ghfs/src/serverError"
	"mjpclab.dev/ghfs/src/setting"
	"mjpclab.dev/ghfs/src/version"
	"os"
	"os/signal"
	"syscall"
)

func cleanupOnEnd(appInst *app.App) {
	chSignal := make(chan os.Signal)
	signal.Notify(chSignal, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-chSignal
		appInst.Close()
		os.Exit(0)
	}()
}

func reopenLogOnHup(appInst *app.App) {
	chSignal := make(chan os.Signal)
	signal.Notify(chSignal, syscall.SIGHUP)

	go func() {
		for range chSignal {
			errs := appInst.ReOpenLog()
			serverError.CheckFatal(errs...)
		}
	}()
}

func main() {
	// params
	params, printVersion, printHelp, errs := param.ParseFromCli()
	serverError.CheckFatal(errs...)
	if printVersion {
		version.PrintVersion()
		os.Exit(0)
	}
	if printHelp {
		param.PrintHelp()
		os.Exit(0)
	}

	// setting
	setting := setting.ParseFromEnv()

	// app
	appInst, errs := app.NewApp(params, setting)
	serverError.CheckFatal(errs...)
	if appInst == nil {
		serverError.CheckFatal(errors.New("failed to create application instance"))
	}

	cleanupOnEnd(appInst)
	reopenLogOnHup(appInst)
	errs = appInst.Open()
	serverError.CheckFatal(errs...)
	appInst.Close()
}
