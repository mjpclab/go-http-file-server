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
		appInst.Shutdown()
	}()
}

func reInitOnHup(appInst *app.App) {
	chSignal := make(chan os.Signal)
	signal.Notify(chSignal, syscall.SIGHUP)

	go func() {
		for range chSignal {
			errs := appInst.ReOpenLog()
			if serverError.CheckError(errs...) {
				appInst.Shutdown()
				break
			}
			errs = appInst.ReLoadCertificates()
			if serverError.CheckError(errs...) {
				appInst.Shutdown()
				break
			}
		}
	}()
}
