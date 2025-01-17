package setting

import (
	"mjpclab.dev/ghfs/src/util"
	"os"
)

type Setting struct {
	CPUProfileFile string
	PidFile        string
	LogQueueSize   string
	Quiet          bool
}

func ParseFromEnv() *Setting {
	cpuProfileFile := os.Getenv("GHFS_CPU_PROFILE_FILE")
	pidFile := os.Getenv("GHFS_PID_FILE")
	logQueueSize := os.Getenv("GHFS_LOG_QUEUE_SIZE")
	quiet := util.GetBoolEnv("GHFS_QUIET")

	return &Setting{
		CPUProfileFile: cpuProfileFile,
		PidFile:        pidFile,
		LogQueueSize:   logQueueSize,
		Quiet:          quiet,
	}
}
