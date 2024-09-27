package middleware

import (
	"mjpclab.dev/ghfs/src/serverLog"
	"os"
)

type Context struct {
	PrefixReqPath string
	VhostReqPath  string
	AliasReqPath  string
	AliasFsPath   string
	AliasFsRoot   string

	WantJson bool

	RestrictAccess bool
	AllowAccess    bool

	NeedAuth     bool
	AuthUserName string
	AuthSuccess  bool

	Status *int

	File     **os.File
	FileInfo *os.FileInfo

	Logger *serverLog.Logger
}
