package setting

import (
	"mjpclab.dev/ghfs/src/util"
	"os"
)

type Setting struct {
	Quiet   bool
	PidFile string
}

func ParseFromEnv() *Setting {
	quiet := util.GetBoolEnv("GHFS_QUIET")
	pidFile := os.Getenv("GHFS_PID_FILE")

	return &Setting{
		Quiet:   quiet,
		PidFile: pidFile,
	}
}
