package setting

import (
	"mjpclab.dev/ghfs/src/util"
	"os"
)

type Setting struct {
	PidFile string
	Quiet   bool
}

func ParseFromEnv() *Setting {
	pidFile := os.Getenv("GHFS_PID_FILE")
	quiet := util.GetBoolEnv("GHFS_QUIET")

	return &Setting{
		PidFile: pidFile,
		Quiet:   quiet,
	}
}
