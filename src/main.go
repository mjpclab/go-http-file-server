package src

import (
	"errors"
	"mjpclab.dev/ghfs/src/app"
	"mjpclab.dev/ghfs/src/param"
	"mjpclab.dev/ghfs/src/serverError"
	"mjpclab.dev/ghfs/src/serverLog"
	"mjpclab.dev/ghfs/src/setting"
	"mjpclab.dev/ghfs/src/version"
	"strconv"
)

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

	// start
	errs = Start(settings, params)
	if serverError.CheckError(errs...) {
		return
	}

	return true
}

func Start(settings *setting.Setting, params param.Params) (errs []error) {
	// CPU profile
	if len(settings.CPUProfileFile) > 0 {
		cpuProfileFile, err := startCPUProfile(settings.CPUProfileFile)
		if err != nil {
			return []error{err}
		}
		defer stopCPUProfile(cpuProfileFile)
	}

	// pid file
	if len(settings.PidFile) > 0 {
		errs = writePidFile(settings.PidFile)
		if len(errs) > 0 {
			return
		}
	}

	// log queue size
	if len(settings.LogQueueSize) > 0 {
		logQueueSize, err := strconv.Atoi(settings.LogQueueSize)
		if err == nil && logQueueSize > 0 {
			serverLog.SetLogQueueSize(logQueueSize)
		}
	}

	// app
	appInst, errs := app.NewApp(params)
	if len(errs) > 0 {
		return
	}
	if appInst == nil {
		errs = []error{errors.New("failed to create application instance")}
		return
	}

	cleanupOnEnd(appInst)
	reInitOnHup(appInst)
	if !settings.Quiet {
		printAccessibleURLs(appInst.GetAccessibleUrls(false))
	}
	errs = appInst.Open()
	return
}
