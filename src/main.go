package src

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

func Main() (ok bool) {
	// params
	params, printVersion, printHelp, errs := param.ParseFromCli()
	if serverError.CheckError(errs...) {
		return
	}
	if printVersion {
		version.PrintVersion()
		return true
	}
	if printHelp {
		param.PrintHelp()
		return true
	}

	// settings
	settings := setting.ParseFromEnv()

	// CPU profile
	if len(settings.CPUProfileFile) > 0 {
		cpuProfileFile, err := StartCPUProfile(settings.CPUProfileFile)
		if serverError.CheckError(err) {
			return
		}
		defer StopCPUProfile(cpuProfileFile)
	}

	// app
	appInst, errs := app.NewApp(params, settings)
	if serverError.CheckError(errs...) {
		return
	}
	if appInst == nil {
		serverError.CheckError(errors.New("failed to create application instance"))
		return
	}

	cleanupOnEnd(appInst)
	reInitOnHup(appInst)
	errs = appInst.Open()
	if serverError.CheckError(errs...) {
		return
	}

	return true
}
