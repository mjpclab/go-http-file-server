package middleware

type Context struct {
	PrefixReqPath string
	VhostReqPath  string
	AliasReqPath  string
	AliasFsPath   string
	AliasFsRoot   string
}
