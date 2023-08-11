package middleware

import (
	"mjpclab.dev/ghfs/src/serverLog"
	"mjpclab.dev/ghfs/src/user"
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

	CanUpload  *bool
	CanMkdir   *bool
	CanDelete  *bool
	CanArchive *bool

	Status *int

	File     **os.File
	FileInfo *os.FileInfo

	Users  *user.List
	Logger *serverLog.Logger
}
