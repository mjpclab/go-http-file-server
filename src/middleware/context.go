package middleware

import "os"

type Context struct {
	PrefixReqPath string
	VhostReqPath  string
	AliasReqPath  string
	AliasFsPath   string
	AliasFsRoot   string

	Item     os.FileInfo
	SubItems []os.FileInfo

	Status int
}
