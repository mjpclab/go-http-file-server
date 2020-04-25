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

	DefaultSort string
	DirIndexes  []string
	Aliases     map[string]string

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

	Shows     *regexp.Regexp
	ShowDirs  *regexp.Regexp
	ShowFiles *regexp.Regexp
	Hides     *regexp.Regexp
	HideDirs  *regexp.Regexp
	HideFiles *regexp.Regexp

	AccessLog string
	ErrorLog  string
}

func (param *Param) GetDupUserNames() []string {
	usersGroups := [][]*user{
		param.UsersPlain,
		param.UsersBase64,
		param.UsersMd5,
		param.UsersSha1,
		param.UsersSha256,
		param.UsersSha512,
	}

	userMap := map[string]bool{}
	dupUserMap := map[string]bool{}

	for _, users := range usersGroups {
		for _, user := range users {
			if userMap[user.Username] {
				dupUserMap[user.Username] = true
			}
			userMap[user.Username] = true
		}
	}

	dupUserNames := make([]string, 0, len(dupUserMap))
	for username, _ := range dupUserMap {
		dupUserNames = append(dupUserNames, username)
	}
	return dupUserNames
}
