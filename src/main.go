package main

import (
	"./app"
	"./param"
	"./serverError"
	"./version"
	"errors"
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

	appInst, errs := app.NewApp(params)
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
