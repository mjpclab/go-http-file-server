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

	NeedAuth     bool
	AuthUserName string
	AuthSuccess  bool

	RestrictAccess bool
	AllowAccess    bool

	WantJson bool

	Status int

	Item     os.FileInfo
	SubItems []os.FileInfo

	Logger *serverLog.Logger
}
