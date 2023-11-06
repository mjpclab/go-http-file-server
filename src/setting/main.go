package setting

import (
	"mjpclab.dev/ghfs/src/util"
	"os"
)

type Setting struct {
	CPUProfileFile string
	PidFile        string
	Quiet          bool
}

func ParseFromEnv() *Setting {
	cpuProfileFile := os.Getenv("GHFS_CPU_PROFILE_FILE")
	pidFile := os.Getenv("GHFS_PID_FILE")
	quiet := util.GetBoolEnv("GHFS_QUIET")

	return &Setting{
		CPUProfileFile: cpuProfileFile,
		PidFile:        pidFile,
		Quiet:          quiet,
	}
}
