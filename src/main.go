package main

import (
	"./app"
	"./param"
	"./serverErrHandler"
	"errors"
	"os"
	"os/signal"
	"syscall"
)

func cleanupOnInterrupt(appInst *app.App) {
	chSignal := make(chan os.Signal)
	signal.Notify(chSignal, syscall.SIGINT)

	go func() {
		<-chSignal
		appInst.Close()
		os.Exit(0)
	}()
}

func reOpenLogOnHup(appInst *app.App) {
	chSignal := make(chan os.Signal)
	signal.Notify(chSignal, syscall.SIGHUP)

	go func() {
		for range chSignal {
			errs := appInst.ReOpenLog()
			serverErrHandler.CheckFatal(errs...)
		}
	}()
}

func main() {
	params := param.ParseCli()
	appInst, errs := app.NewApp(params)
	serverErrHandler.CheckFatal(errs...)

	if appInst == nil {
		serverErrHandler.CheckFatal(errors.New("failed to create application instance"))
	}

	cleanupOnInterrupt(appInst)
	reOpenLogOnHup(appInst)
	errs = appInst.Open()
	serverErrHandler.CheckFatal(errs...)

	appInst.Close()
}
