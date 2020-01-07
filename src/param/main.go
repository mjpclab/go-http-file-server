package param

import "regexp"

type user struct {
	Username string
	Password string
}

type Param struct {
	Root      string
	EmptyRoot bool

	FallbackProxies          map[string]string
	AlwaysProxies            map[string]string
	IgnoreProxyTargetBadCert bool

	Aliases map[string]string

	GlobalUpload bool
	UploadUrls   []string
	UploadDirs   []string

	GlobalArchive bool
	ArchiveUrls   []string
	ArchiveDirs   []string

	GlobalCors bool
	CorsUrls   []string
	CorsDirs   []string

	GlobalAuth  bool
	AuthUrls    []string
	AuthDirs    []string
	UsersPlain  []*user
	UsersBase64 []*user
	UsersMd5    []*user
	UsersSha1   []*user
	UsersSha256 []*user
	UsersSha512 []*user

	Key          string
	Cert         string
	Listens      []string
	ListensPlain []string
	ListensTLS   []string
	HostNames    []string
	Template     string

	Shows     *regexp.Regexp
	ShowDirs  *regexp.Regexp
	ShowFiles *regexp.Regexp
	Hides     *regexp.Regexp
	HideDirs  *regexp.Regexp
	HideFiles *regexp.Regexp

	AccessLog string
	ErrorLog  string
}
