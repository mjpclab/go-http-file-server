package param

import (
	"crypto/tls"
	"regexp"
)

type user struct {
	Username string
	Password string
}

type Param struct {
	Root      string
	EmptyRoot bool

	DefaultSort   string
	DirIndexes    []string
	Aliases       map[string]string
	GlobalHeaders [][2]string

	GlobalUpload bool
	UploadUrls   []string
	UploadDirs   []string

	GlobalMkdir bool
	MkdirUrls   []string
	MkdirDirs   []string

	GlobalDelete bool
	DeleteUrls   []string
	DeleteDirs   []string

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

	Certificate  *tls.Certificate
	Listens      []string
	ListensPlain []string
	ListensTLS   []string
	HostNames    []string
	Template     string

	GlobalHsts  bool
	GlobalHttps bool
	HttpsPort   string

	Shows     *regexp.Regexp
	ShowDirs  *regexp.Regexp
	ShowFiles *regexp.Regexp
	Hides     *regexp.Regexp
	HideDirs  *regexp.Regexp
	HideFiles *regexp.Regexp

	AccessLog string
	ErrorLog  string
}
