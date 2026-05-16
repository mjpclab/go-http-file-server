package src

import (
	"os"
	"os/signal"
	"syscall"

	"mjpclab.dev/ghfs/src/app"
	"mjpclab.dev/ghfs/src/serverError"
)

func cleanupOnEnd(appInst *app.App) {
	chSignal := make(chan os.Signal)
	signal.Notify(chSignal, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-chSignal
		appInst.Close()
	}()
}

func reInitOnHup(appInst *app.App) {
	chSignal := make(chan os.Signal)
	signal.Notify(chSignal, syscall.SIGHUP)

	go func() {
		for _ = range chSignal {
			errs := appInst.ReOpenLog()
			if serverError.CheckError(errs...) {
				appInst.Close()
				break
			}
			errs = appInst.ReLoadCertificates()
			if serverError.CheckError(errs...) {
				appInst.Close()
				break
			}
		}
	}()
}
