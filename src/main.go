package main

import (
	"./app"
	"./param"
	"os"
	"os/signal"
	"syscall"
)

var appInst *app.App

func cleanupOnInterrupt() {
	chSignal := make(chan os.Signal)
	signal.Notify(chSignal, syscall.SIGINT)

	go func() {
		<-chSignal
		appInst.Close()
		os.Exit(0)
	}()
}

func reOpenLogOnHup() {
	chSignal := make(chan os.Signal)
	signal.Notify(chSignal, syscall.SIGHUP)

	go func() {
		for range chSignal {
			appInst.ReOpenLog()
		}
	}()
}

func main() {
	cleanupOnInterrupt()

	params := param.ParseCli()
	appInst = app.NewApp(params)

	if appInst != nil {
		reOpenLogOnHup()
		appInst.Open()
		defer appInst.Close()
	}
}
